# Prod uses a separate Algolia application (different ALGOLIA_APPLICATION_ID/API_KEY)
# for data isolation; index names match dev.
resource "algolia_index" "users_index" {
  name = "users"

  searchable_attributes = [
    "username",
  ]
}
