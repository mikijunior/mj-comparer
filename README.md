# mj-comparer
The project is helpful for comare images.
## Algorithms
1. Compare by pixels. Algorithm will be used for images with the same extensions. Here we would be able to compare every pixels of images to check their color codes.
2. Compare by Euclidean distance. Algorithm will be used for images with different extensions. Images are resized in a special way to squares of fixed size called "icons". Euclidean distance between the icons is used to give the similarity verdict. Also image proportions are used to avoid matching images of distinct shape. Read more:  https://vitali-fedulov.github.io/similar.pictures/algorithm-for-perceptual-image-comparison.html

## How to run the project
- Create new config file by example in <project_path>/configs and fill it for your configurations.
- Run migration with go-migrate. Example: `migrate -path ./migrations/ -database "mysql://<user>:<password>@tcp(<DB_URL>)/<DB_NAME>" -verbose up`
### Request
Header:
```JSON
{
  "Api-key": "Key from db"
}
```
Body:
```JSON
{
  "src": "Link to image",
  "dst": "Link to image",
}
```
### Response
```JSON
{
    "Is_same": true,
    "Percent": 0,
    "Algorithm": "Euclidean distance / Pixel compare"
}
```
