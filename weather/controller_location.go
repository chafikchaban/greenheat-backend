package weather

import (
	"encoding/json"
	"fmt"
	"strings"
)

// LocationController is Controller that handles operations on locations
type LocationController struct {
}

// GenerateID generates a unique identifier for the location based on latitude and longitude
func GenerateID(latitude, longitude string) string {
    return strings.Join([]string{latitude, longitude}, "_")
}

// AddLocation adds a new location to the database if it's unique
func (lc *LocationController) AddLocation(db Database, newLocation Location) error {
    // Generate a unique ID from the coordinates
    newID := GenerateID(newLocation.Latitude, newLocation.Longitude)
    newLocation.ID = newID

    // Check if the location already exists
    var existing Location
    err := db.d.Read("locations", newID, &existing)
    if err == nil {
        return fmt.Errorf("location with latitude %s and longitude %s already exists", newLocation.Latitude, newLocation.Longitude)
    }

    // Save the new location to the database
    if err := db.d.Write("locations", newID, newLocation); err != nil {
        return fmt.Errorf("could not save location: %v", err)
    }

    return nil
}

// GetLocations retrieves all locations from the database
func (lc *LocationController) GetLocations(db Database) ([]Location, error) {
    var locations []Location

    // Get all records from the "locations" collection
    records, err := db.d.ReadAll("locations")
    if err != nil {
        return nil, fmt.Errorf("could not read locations: %v", err)
    }

    // Unmarshal each record into the locations slice
    for _, record := range records {
        var location Location
        if err := json.Unmarshal([]byte(record), &location); err != nil {
            return nil, fmt.Errorf("could not unmarshal location: %v", err)
        }
        locations = append(locations, location)
    }

    return locations, nil
}

// DeleteLocation removes a location from the database based on its unique ID
func (lc *LocationController) DeleteLocation(db Database, id string) error {

    // Check if the location exists in the database
    var location Location
    err := db.d.Read("locations", id, &location)
    if err != nil {
        return fmt.Errorf("location with id %s does not exist", id)
    }

    // Delete the location from the database
    if err := db.d.Delete("locations", id); err != nil {
        return fmt.Errorf("could not delete location with ID %s: %v", id, err)
    }

    fmt.Printf("Location %s with ID %s deleted successfully.\n", location.Name, id)
    return nil
}

// InitializeLocations populates the database with initial location data for each German state
func (lc *LocationController) InitializeLocations(db Database) {
    // List of coordinates for each German state
    locations := []Location{
        {Name: "Baden-Württemberg", Latitude: "48.6616", Longitude: "9.3501"},
        {Name: "Bavaria", Latitude: "48.7904", Longitude: "11.4979"},
        {Name: "Berlin", Latitude: "52.5200", Longitude: "13.4050"},
        {Name: "Brandenburg", Latitude: "52.4125", Longitude: "12.5316"},
        {Name: "Bremen", Latitude: "53.0793", Longitude: "8.8017"},
        {Name: "Hamburg", Latitude: "53.5511", Longitude: "9.9937"},
        {Name: "Hesse", Latitude: "50.6521", Longitude: "9.1624"},
        {Name: "Lower Saxony", Latitude: "52.6367", Longitude: "9.8451"},
        {Name: "Mecklenburg-Vorpommern", Latitude: "53.6127", Longitude: "12.4296"},
        {Name: "North Rhine-Westphalia", Latitude: "51.4332", Longitude: "7.6616"},
        {Name: "Rhineland-Palatinate", Latitude: "49.9454", Longitude: "7.4514"},
        {Name: "Saarland", Latitude: "49.3964", Longitude: "7.0236"},
        {Name: "Saxony", Latitude: "51.1045", Longitude: "13.2017"},
        {Name: "Saxony-Anhalt", Latitude: "51.9506", Longitude: "11.6928"},
        {Name: "Schleswig-Holstein", Latitude: "54.2194", Longitude: "9.6961"},
        {Name: "Thuringia", Latitude: "51.0101", Longitude: "11.1637"},
    }

    // Iterate through each location and add it to the database if not present
    for _, location := range locations {
        // Generate the unique ID from the coordinates
        location.ID = GenerateID(location.Latitude, location.Longitude)

        // Check if the location already exists
        var existing Location
        err := db.d.Read("locations", location.ID, &existing)
        if err == nil {
            fmt.Printf("Location %s with ID %s already exists. Skipping...\n", location.Name, location.ID)
            continue // If location already exists, skip adding it
        }

        // Save the location to the database
        if err := db.d.Write("locations", location.ID, location); err != nil {
            fmt.Printf("Error saving location %s: %v\n", location.Name, err)
        } else {
            fmt.Printf("Location %s saved successfully.\n", location.Name)
        }
    }
}

