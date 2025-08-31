# Query all users from the BMLT server
data "bmlt_users" "all" {}

# Output the list of users
output "all_users" {
  description = "All users in the BMLT server"
  value = {
    count = length(data.bmlt_users.all.users)
    users = [
      for user in data.bmlt_users.all.users : {
        id           = user.id
        username     = user.username
        display_name = user.display_name
        type         = user.type
      }
    ]
  }
}
