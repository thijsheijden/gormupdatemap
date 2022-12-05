# gormupdatemap

## Introduction
A GORM helper package for easy creation of maps to update DB records with.

When updating a struct in GORM using the `Updates` method, default values are ignored. This makes sense as otherwise all fields of a struct would have to be set to prevent the resetting of field values. However, if we want to set default values, like `false`, `0` and `""` we need to either specify exact columns to update, or use a map.

This package adds the `CreateUpdateMap` function, which allows you to easily create such a update map to use when updating records.

## Installation
`go get github.com/thijsheijden/gormupdatemap`

## Usage and example
### Update struct type
To use the `CreateUpdateMap` function you first need to define a struct which contains the fields that can be updated.

Let's say we have the following `Person` database model:

```go
type Person struct {
  GivenName   *string `json:"given_name"`
  FamilyName  *string `json:"family_name"`
  Age         *int    `json:"age"`
  Verified    *bool   `json:"verified"`
}
```

If we want to create an update type for this model, that allows the updating of all fields, we could create the following type:

```go
type PersonUpdate struct {
  GivenName   *string `json:"given_name"`
  FamilyName  *string `json:"family_name"`
  Age         *int    `json:"age"`
  Verified    *bool   `json:"verified"`
}
```

### Pointers
It is important to note that all fields are pointers. This is used to determine if a field should be updated or not. Nil values are ignored when creating the update map.

### JSON tags
The JSON tags are also important, as they are used as the column names when updating (the keys in the update map). A very important detail here is that GORM will modify column names containing numbers. For instance the field `Field1` would have the column name `field1` and not the expected `field_1`. The `CreateUpdateMap` function takes this into account and will convert JSON tags with numbers into the valid GORM column names.

If we wanted to only allow updating of the `Verified` field, we would use the following type:

```go
type PersonUpdate struct {
  Verified *bool `json:"verified"`
}
```

### Admin only fields
Sometimes we want certain fields of a model to be updatable only by admins. This can be set by adding the `admin_only` tag to those fields like this:

```go
type PersonUpdate struct {
  Verified *bool `json:"verified" admin_only:"true"`
}
```

### Parameters and return values
`CreateUpdateMap` takes in two parameters. The first is a struct with pointer fields describing columns to be updated. The second is a boolean denoting whether the caller is an admin or not.

Two values are returned. The first value is the map that can be used to perform the DB update with. The second value is a pointer to a string, which will contain any validation error that occurred. For now this will only be attempting to update an admin-only field as a non-admin.

### Example
```go
newName := "John"
newFamilyName := "Doe"

personUpdate := PersonUpdate{
  GivenName: &newName,
  FamilyName: &newFamilyName,
}

updateMap, validationErr := gormupdatemap.CreateUpdateMap(personUpdate, false)
if validationErr != nil {
  // HTTP Bad Request, or Unauthorized
}

// db is a GORM DB
// person is the Person model we want to update
if err := db.Model(&person).Updates(updateMap).Error; err != nil {
  // Internal server error
}
```