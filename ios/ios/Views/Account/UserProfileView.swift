//
//  UserProfileView.swift
//  ios
//

import PhotosUI
import SwiftUI
import UIKit

struct UserProfileView: View {
  @Environment(AuthenticationManager.self) private var authManager
  @Environment(EventReporterService.self) private var eventReporterService
  @Environment(\.dismiss) private var dismiss
  @State private var profileViewModel: UserProfileViewModel?
  @State private var avatarViewModel: UploadAvatarViewModel?
  @State private var selectedAvatarItem: PhotosPickerItem?
  @State private var showDetailsConfirmSheet = false
  @State private var confirmPassword = ""
  @State private var confirmTotpToken = ""

  var body: some View {
    Group {
      if let profileViewModel = profileViewModel {
        ScrollView {
          VStack(spacing: DSTheme.Spacing.xl) {
            profilePhotoSection
            userDetailsSection(viewModel: profileViewModel)
          }
          .dsScreenPadding()
        }
      } else {
        DSInitializingView()
      }
    }
    .navigationTitle("Profile")
    .refreshable {
      await profileViewModel?.loadUser()
    }
    .onAppear {
      if profileViewModel == nil {
        profileViewModel = UserProfileViewModel(authManager: authManager)
        eventReporterService.reporter.track(event: "profile_viewed", properties: [:])
      }
      if avatarViewModel == nil {
        avatarViewModel = UploadAvatarViewModel(authManager: authManager)
      }
      Task {
        await profileViewModel?.loadUser()
      }
    }
    .onChange(of: selectedAvatarItem) { _, newItem in
      guard let item = newItem, let avatarViewModel = avatarViewModel else { return }
      eventReporterService.reporter.track(event: "profile_avatar_selected", properties: [:])
      Task {
        await loadAndUploadAvatar(item: item, viewModel: avatarViewModel)
      }
    }
    .onChange(of: avatarViewModel?.didSucceed) { _, didSucceed in
      if didSucceed == true {
        eventReporterService.reporter.track(event: "profile_avatar_uploaded", properties: [:])
        Task { await profileViewModel?.loadUser() }
      }
    }
    .sheet(isPresented: $showDetailsConfirmSheet) {
      confirmDetailsSheet
    }
  }

  // MARK: - Profile Photo Section

  private var profilePhotoSection: some View {
    DSSection("Profile Photo", subtitle: "Select a photo to set as your profile picture") {
      if let avatarViewModel = avatarViewModel {
        VStack(spacing: DSTheme.Spacing.lg) {
          PhotosPicker(
            selection: $selectedAvatarItem,
            matching: .images,
            photoLibrary: .shared()
          ) {
            HStack(spacing: DSTheme.Spacing.md) {
              Image(systemName: "person.crop.circle.badge.plus")
                .font(.system(size: DSTheme.IconSize.lg))
                .foregroundColor(DSTheme.Colors.primary)
              Text("Select Photo")
                .font(DSTheme.Typography.label)
                .foregroundColor(DSTheme.Colors.textPrimary)
            }
            .frame(maxWidth: .infinity)
            .padding(DSTheme.Spacing.lg)
            .background(DSTheme.Colors.cardBackground)
            .cornerRadius(DSTheme.Radius.md)
          }
          .disabled(avatarViewModel.isUploading)

          if avatarViewModel.isUploading {
            ProgressView("Uploading...")
              .padding()
          }

          if let error = avatarViewModel.errorMessage {
            Text(error)
              .font(DSTheme.Typography.caption)
              .foregroundColor(DSTheme.Colors.error)
              .multilineTextAlignment(.leading)
          }

          if avatarViewModel.didSucceed {
            HStack(spacing: DSTheme.Spacing.sm) {
              Image(systemName: "checkmark.circle.fill")
                .foregroundColor(DSTheme.Colors.success)
              Text("Profile photo updated successfully")
                .font(DSTheme.Typography.label)
                .foregroundColor(DSTheme.Colors.success)
            }
            .frame(maxWidth: .infinity, alignment: .leading)
            .padding(DSTheme.Spacing.md)
            .background(DSTheme.Colors.cardBackground)
            .cornerRadius(DSTheme.Radius.md)
          }
        }
      }
    }
  }

  // MARK: - User Details Section

  private func userDetailsSection(viewModel: UserProfileViewModel) -> some View {
    DSSection("User Details", subtitle: "Update your profile information") {
      VStack(spacing: DSTheme.Spacing.lg) {
        if let error = viewModel.errorMessage {
          Text(error)
            .font(DSTheme.Typography.caption)
            .foregroundColor(DSTheme.Colors.error)
            .multilineTextAlignment(.leading)
            .frame(maxWidth: .infinity, alignment: .leading)
        }

        DSTextField(
          "Username",
          text: Binding(get: { viewModel.username }, set: { viewModel.username = $0 }),
          type: .username
        )

        if viewModel.usernameHasChanged {
          DSButton(
            "Update Username",
            icon: "checkmark",
            fullWidth: true,
            isDisabled: viewModel.username.isEmpty || viewModel.isLoading
          ) {
            Task {
              await viewModel.updateUsername()
            }
          }
        }

        DSTextField(
          "First Name",
          text: Binding(get: { viewModel.firstName }, set: { viewModel.firstName = $0 })
        )

        DSTextField(
          "Last Name",
          text: Binding(get: { viewModel.lastName }, set: { viewModel.lastName = $0 })
        )

        DatePicker(
          "Birthday",
          selection: Binding(get: { viewModel.birthday }, set: { viewModel.birthday = $0 }),
          displayedComponents: .date
        )
        .font(DSTheme.Typography.body)
        .foregroundColor(DSTheme.Colors.textPrimary)

        DSButton(
          "Update Details",
          icon: "checkmark",
          fullWidth: true,
          isDisabled: !viewModel.detailsHasChanged || viewModel.isLoading
        ) {
          confirmPassword = ""
          confirmTotpToken = ""
          showDetailsConfirmSheet = true
        }
      }
    }
  }

  // MARK: - Confirm Details Sheet

  private var confirmDetailsSheet: some View {
    NavigationStack {
      Form {
        Section {
          Text("Confirm your password to update your profile details.")
            .font(DSTheme.Typography.body)
            .foregroundColor(DSTheme.Colors.textSecondary)
        }
        Section {
          SecureField("Password", text: $confirmPassword)
          if profileViewModel?.hasTwoFactor == true {
            TextField("Authenticator Code", text: $confirmTotpToken)
              .keyboardType(.numberPad)
          }
        }
      }
      .navigationTitle("Confirm Update")
      .navigationBarTitleDisplayMode(.inline)
      .toolbar {
        ToolbarItem(placement: .cancellationAction) {
          Button("Cancel") {
            showDetailsConfirmSheet = false
          }
        }
        ToolbarItem(placement: .confirmationAction) {
          Button("Update") {
            Task {
              if await profileViewModel?.updateUserDetails(
                currentPassword: confirmPassword,
                totpToken: confirmTotpToken
              ) == true {
                eventReporterService.reporter.track(
                  event: "profile_details_updated", properties: [:])
                showDetailsConfirmSheet = false
              }
            }
          }
          .disabled(confirmPassword.isEmpty)
        }
      }
    }
  }

  // MARK: - Avatar Upload

  private func loadAndUploadAvatar(item: PhotosPickerItem, viewModel: UploadAvatarViewModel) async {
    do {
      guard let data = try await item.loadTransferable(type: Data.self) else {
        viewModel.errorMessage = "Failed to load image data"
        return
      }

      let prepared = ImagePreparationHelper.prepareImageForUpload(data: data)
      let objectName = UUID().uuidString + "." + prepared.fileExtension

      await viewModel.uploadAvatar(
        imageData: prepared.data,
        contentType: prepared.contentType,
        objectName: objectName
      )

      selectedAvatarItem = nil
    } catch {
      viewModel.errorMessage = "Failed to load image: \(error.localizedDescription)"
    }
  }
}

#Preview {
  let authManager = AuthenticationManager()
  authManager.isAuthenticated = true
  authManager.username = "Test User"
  authManager.userID = "user123"

  return NavigationStack {
    UserProfileView()
      .environment(authManager)
  }
}
