# Image Processor #

This README it's about Image Processor

### Project ###

* GoLang 1.16.4
* Alpine 3.13

### Require ###
* Docker-ce 17.9+
* Docker-compose +1.5

---

### Usage ###

This service allows to do few changes on image on the fly, changes like: 
- Convert to Webp
- Generate the image in grayscale
- Resize the image
- Generate the image with sepia filter
- Generate the image with blur filter
- Flip image vertically 
- Flip image horizontally
- Compression quality

### Convert Image to Webp ###

- Add parameter `webp=1` as query string

### Generate Image in Grayscale ###

- Add parameter `grayscale=true` as query string


### Resize Image ###

- Add parameter `width=100` and `height=100` as query string


### Generate Image with Sepia filter ###

- Add parameter `sepia=true` as query string


### Generate Image with blur filter ###

- Add parameter `blur=10` (the blur level is 1 ~ 100) as query string


### Flip image vertically ###

- Add parameter `flipv=true` as query string


### Flip image horizontally ###

- Add parameter `fliph=true` as query string


### Change compression quality ###

- Add parameter `quality=75` (the quality legel is 0 ~ 100) as query string

### Crop image ###

- Add parameter `crop=100x100x200x300`. The values as parameters are mandatory when you are using crop and should the values are `InitialHorizontalPoint x InitialVerticalPoint x FinalHorizontalPoint x FinalVerticalPoint`

---

## Initial Settings ##

Copy env.sample to .env and fill these information:
```
IMAGE_BASE_URL=<BASE_DOMAIN_WITH_IMAGES_TO_BE_USED>
IMAGE_PROCESSOR_DOMAIN=<SERVICE_DOMAIN>
IMAGE_PROCESSOR_IP=<SERVICE_IP>
```

### Up Container ###
```
$ sh run-local.sh
```

### Do Request ###
On browser call your domain with this container as you .env file with parameter to load image, ex: 
```
<IMAGE_PROCESSOR_DOMAIN>/bros.jpg?quality=50&webp=1
```

### Sample page with few features running
On browser goes to /doc/index.html and you can see the images converted to webp and with few features running