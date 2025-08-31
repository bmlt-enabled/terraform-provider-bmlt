data "bmlt_formats" "all" {}

# Access format information
output "format_count" {
  value = length(data.bmlt_formats.all.formats)
}

output "format_names" {
  value = [for format in data.bmlt_formats.all.formats : 
    format.translations[0].name if length(format.translations) > 0
  ]
}
