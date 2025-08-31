resource "bmlt_format" "example" {
  world_id = "CUSTOM_FORMAT"
  type     = "FC3"

  translations {
    key         = "en"
    name        = "Custom Format"
    description = "This is a custom format for our region"
    language    = "en"
  }

  translations {
    key         = "es"
    name        = "Formato Personalizado"
    description = "Este es un formato personalizado para nuestra región"
    language    = "es"
  }
}
