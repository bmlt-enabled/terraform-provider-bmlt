---
page_title: "bmlt_meetings Data Source - terraform-provider-bmlt"
subcategory: ""
description: |-
  Meetings data source allows you to retrieve information about meetings.
---

# bmlt_meetings (Data Source)

Meetings data source allows you to retrieve information about meetings.

This data source supports various filtering options to help you find specific meetings based on your criteria.

## Example Usage

### Basic Usage - Get All Meetings

```terraform
# Get all meetings
data "bmlt_meetings" "all" {}

# Get meetings filtered by service body
data "bmlt_meetings" "service_body_meetings" {
  service_body_ids = "1,2,3"
}

# Get meetings for specific days
data "bmlt_meetings" "weekend_meetings" {
  days = "0,6" # Sunday and Saturday
}

# Search meetings by name
data "bmlt_meetings" "searched_meetings" {
  search_string = "Big Book"
}

# Output examples
output "total_meetings" {
  value = length(data.bmlt_meetings.all.meetings)
}

output "weekend_meeting_names" {
  value = [for meeting in data.bmlt_meetings.weekend_meetings.meetings : meeting.name]
}
```

### Filter by Service Body

```hcl
# Get meetings for specific service bodies
data "bmlt_meetings" "area_meetings" {
  service_body_ids = "1,2,3"
}

output "area_meeting_count" {
  value = length(data.bmlt_meetings.area_meetings.meetings)
}
```

### Filter by Day of Week

```hcl
# Get weekend meetings (Saturday and Sunday)
data "bmlt_meetings" "weekend_meetings" {
  days = "0,6" # 0=Sunday, 6=Saturday
}

# Get weekday meetings (Monday through Friday)
data "bmlt_meetings" "weekday_meetings" {
  days = "1,2,3,4,5"
}

output "weekend_meeting_names" {
  value = [for meeting in data.bmlt_meetings.weekend_meetings.meetings : meeting.name]
}
```

### Filter by Multiple Criteria

```hcl
# Get Monday meetings for specific service bodies
data "bmlt_meetings" "monday_area_meetings" {
  service_body_ids = "1,2"
  days            = "1" # Monday only
}

# Search for meetings with specific text
data "bmlt_meetings" "beginners_meetings" {
  search_string = "beginners"
}
```

### Filter by Specific Meeting IDs

```hcl
# Get specific meetings by their IDs
data "bmlt_meetings" "specific_meetings" {
  meeting_ids = "100,200,300"
}

# Use the filtered results
output "specific_meeting_details" {
  value = {
    for meeting in data.bmlt_meetings.specific_meetings.meetings :
    meeting.id => {
      name      = meeting.name
      day       = meeting.day
      time      = meeting.start_time
      location  = meeting.location_text
    }
  }
}
```

### Using Meeting Data for Other Resources

```hcl
# Get all meetings for analysis
data "bmlt_meetings" "all" {}

# Create a summary of meeting distribution by day
locals {
  meetings_by_day = {
    for day in range(7) : day => [
      for meeting in data.bmlt_meetings.all.meetings : meeting
      if meeting.day == day
    ]
  }
  
  day_names = {
    0 = "Sunday"
    1 = "Monday" 
    2 = "Tuesday"
    3 = "Wednesday"
    4 = "Thursday"
    5 = "Friday"
    6 = "Saturday"
  }
}

# Output meeting statistics
output "meetings_per_day" {
  value = {
    for day, meetings in local.meetings_by_day :
    local.day_names[day] => length(meetings)
  }
}

# Find meetings that need location updates
output "meetings_without_coordinates" {
  value = [
    for meeting in data.bmlt_meetings.all.meetings :
    {
      id       = meeting.id
      name     = meeting.name
      location = meeting.location_text
    }
    if meeting.latitude == 0 && meeting.longitude == 0
  ]
}
```

### Virtual and Hybrid Meetings

```hcl
# Get virtual meetings (venue_type would need to be exposed in the schema)
data "bmlt_meetings" "all_meetings" {}

locals {
  # Filter meetings with virtual meeting links
  virtual_meetings = [
    for meeting in data.bmlt_meetings.all_meetings.meetings :
    meeting if meeting.virtual_meeting_link != ""
  ]
}

output "virtual_meeting_links" {
  value = {
    for meeting in local.virtual_meetings :
    meeting.name => meeting.virtual_meeting_link
  }
}
```

### Meeting Format Analysis

```hcl
# Get all meetings and formats for analysis
data "bmlt_meetings" "all" {}
data "bmlt_formats" "all" {}

# Create a map of format IDs to names for easy lookup
locals {
  format_names = {
    for format in data.bmlt_formats.all.formats :
    format.id => format.translations[0].name
    if length(format.translations) > 0
  }
  
  # Count meetings by format
  format_usage = {
    for format_id, format_name in local.format_names :
    format_name => length([
      for meeting in data.bmlt_meetings.all.meetings :
      meeting if contains(meeting.format_ids, format_id)
    ])
  }
}

output "most_common_formats" {
  value = local.format_usage
}
```

## Day of Week Reference

When using the `days` filter, use these numeric values:

- `0` = Sunday
- `1` = Monday  
- `2` = Tuesday
- `3` = Wednesday
- `4` = Thursday
- `5` = Friday
- `6` = Saturday

## Time Format

Meeting times are returned in 24-hour format (HH:MM), for example:
- `"09:00"` = 9:00 AM
- `"19:30"` = 7:30 PM
- `"12:00"` = 12:00 PM (Noon)

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `days` (String) Comma delimited day ids between 0-6 to filter by
- `meeting_ids` (String) Comma delimited meeting ids to filter by
- `search_string` (String) Search string to filter meetings
- `service_body_ids` (String) Comma delimited service body ids to filter by

### Read-Only

- `id` (String) Placeholder identifier for the data source.
- `meetings` (Attributes List) List of meetings (see [below for nested schema](#nestedatt--meetings))

<a id="nestedatt--meetings"></a>
### Nested Schema for `meetings`

Read-Only:

- `bus_lines` (String) Bus lines
- `comments` (String) Comments
- `contact_email_1` (String) Primary contact email
- `contact_email_2` (String) Secondary contact email
- `contact_name_1` (String) Primary contact name
- `contact_name_2` (String) Secondary contact name
- `contact_phone_1` (String) Primary contact phone
- `contact_phone_2` (String) Secondary contact phone
- `day` (Number) Day of the week (0=Sunday, 1=Monday, etc.)
- `duration` (String) Meeting duration
- `email` (String) Meeting email
- `format_ids` (List of Number) List of format identifiers
- `id` (Number) Meeting identifier
- `latitude` (Number) Latitude coordinate
- `location_city_subsection` (String) City subsection
- `location_info` (String) Location info
- `location_municipality` (String) Municipality
- `location_nation` (String) Nation
- `location_neighborhood` (String) Neighborhood
- `location_postal_code_1` (String) Postal code
- `location_province` (String) Province
- `location_street` (String) Street address
- `location_sub_province` (String) Sub province
- `location_text` (String) Location text
- `longitude` (Number) Longitude coordinate
- `name` (String) Meeting name
- `phone_meeting_number` (String) Phone meeting number
- `published` (Boolean) Whether the meeting is published
- `service_body_id` (Number) Service body identifier
- `start_time` (String) Meeting start time
- `temporarily_virtual` (Boolean) Whether the meeting is temporarily virtual
- `time_zone` (String) Time zone
- `train_lines` (String) Train lines
- `venue_type` (Number) Venue type (1=in-person, 2=virtual, 3=hybrid)
- `virtual_meeting_additional_info` (String) Additional virtual meeting info
- `virtual_meeting_link` (String) Virtual meeting link
- `world_id` (String) World identifier
