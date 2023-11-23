# mj-comparer

The project is helpful for comparing images.

## Algorithms

1. **Structural Similarity Index (SSI):** The SSI algorithm is implemented to provide accurate image comparisons. It's a perception-based model that considers various features of images, such as structural information, luminance, and contrast.

    Read more: [Structural Similarity Index](https://vitali-fedulov.github.io/similar.pictures/algorithm-for-perceptual-image-comparison.html)

## How to run the project

- Create a new config file by following the example in `<project_path>/configs` and fill it with your configurations.
- Run migration with go-migrate: `make migrate-up u=<db_username> p=<db_password> host=<db_host> name=<db_name>`
- Run the app: `make run`
- Build the app: `make build`

### Request

Header:
```json
{
  "Api-key": "Key from db"
}
```
Body:
```json
{
  "src": "Link to image",
  "dst": "Link to image",
}
```
Response: 
```json
{
    "Is_same": true,
    "Percent": 97.34,
    "Algorithm": "Structural Similarity Index"
}
```