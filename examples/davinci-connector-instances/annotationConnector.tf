resource "pingone_davinci_connector_instance" "annotationConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "annotationConnector"
  }
  name = "My awesome annotationConnector"
  property {
    name  = "annotation"
    type  = "string"
    value = var.annotationconnector_property_annotation
  }
  property {
    name  = "annotationTextColor"
    type  = "string"
    value = var.annotationconnector_property_annotation_text_color
  }
  property {
    name  = "backgroundColor"
    type  = "string"
    value = var.annotationconnector_property_background_color
  }
  property {
    name  = "cornerRadius"
    type  = "string"
    value = var.annotationconnector_property_corner_radius
  }
  property {
    name  = "fontFamily"
    type  = "string"
    value = var.annotationconnector_property_font_family
  }
  property {
    name  = "fontSize"
    type  = "string"
    value = var.annotationconnector_property_font_size
  }
  property {
    name  = "fontStyle"
    type  = "string"
    value = var.annotationconnector_property_font_style
  }
  property {
    name  = "height"
    type  = "string"
    value = var.annotationconnector_property_height
  }
  property {
    name  = "strokeColor"
    type  = "string"
    value = var.annotationconnector_property_stroke_color
  }
  property {
    name  = "strokeEnabled"
    type  = "string"
    value = var.annotationconnector_property_stroke_enabled
  }
  property {
    name  = "strokeWidth"
    type  = "string"
    value = var.annotationconnector_property_stroke_width
  }
  property {
    name  = "width"
    type  = "string"
    value = var.annotationconnector_property_width
  }
}
