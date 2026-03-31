//
//  ImagePreparationHelper.swift
//  ios
//

import UIKit

/// Prepares image data for upload. Converts HEIC to JPEG when needed, preserves
/// format for PNG/JPEG/GIF. Use for avatar, recipe, meal, or any image upload.
public enum ImagePreparationHelper {
  public struct PreparedImage: Sendable {
    public let data: Data
    public let contentType: String
    public let fileExtension: String

    public init(data: Data, contentType: String, fileExtension: String) {
      self.data = data
      self.contentType = contentType
      self.fileExtension = fileExtension
    }
  }

  /// Prepares image data for upload. Converts HEIC to JPEG since backend doesn't support HEIC.
  public static func prepareImageForUpload(data: Data) -> PreparedImage {
    if isHEIC(data: data), let image = UIImage(data: data),
      let jpegData = image.jpegData(compressionQuality: 0.9)
    {
      return PreparedImage(data: jpegData, contentType: "image/jpeg", fileExtension: "jpg")
    }
    if isPNG(data: data) {
      return PreparedImage(data: data, contentType: "image/png", fileExtension: "png")
    }
    if isJPEG(data: data) {
      return PreparedImage(data: data, contentType: "image/jpeg", fileExtension: "jpg")
    }
    if isGIF(data: data) {
      return PreparedImage(data: data, contentType: "image/gif", fileExtension: "gif")
    }
    if let image = UIImage(data: data), let jpegData = image.jpegData(compressionQuality: 0.9) {
      return PreparedImage(data: jpegData, contentType: "image/jpeg", fileExtension: "jpg")
    }
    return PreparedImage(data: data, contentType: "image/jpeg", fileExtension: "jpg")
  }

  public static func isJPEG(data: Data) -> Bool {
    data.count >= 3 && data[0] == 0xFF && data[1] == 0xD8 && data[2] == 0xFF
  }

  public static func isPNG(data: Data) -> Bool {
    let signature: [UInt8] = [0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A]
    return data.count >= 8 && data.prefix(8).elementsEqual(signature)
  }

  public static func isGIF(data: Data) -> Bool {
    let signature = "GIF8"
    return data.count >= 4 && String(data: data.prefix(4), encoding: .ascii) == signature
  }

  public static func isHEIC(data: Data) -> Bool {
    guard data.count >= 12 else { return false }
    let ftyp = String(data: data.subdata(in: 4..<8), encoding: .ascii)
    guard ftyp == "ftyp" else { return false }
    let brand = String(data: data.subdata(in: 8..<12), encoding: .ascii)
    return brand == "heic" || brand == "heix" || brand == "mif1"
  }
}
