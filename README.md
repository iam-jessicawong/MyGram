## About MyGram
MyGram is a web service where user can post photo, add comment on other user photos and add social media. This project has 4 group endpoint which is:
1. User :
    - Register [POST] 
    - Login [POST]

2. Photo : 
    - GetAll [GET]
    - GetOne [GET] 
    - CreatePhoto [POST] 
    - UpdatePhoto [PUT] 
    - DeletePhoto [DELETE]

3. Comment :
    - GetAll [GET] 
    - GetOne [GET]
    - CreateComment [POST] 
    - UpdateComment [PUT] 
    - DeleteComment [DELETE]

4. Social Media :
    - GetAll [GET] 
    - GetOne [GET]
    - CreateSocialMedia [POST]
    - UpdateSocialMedia [PUT]
    - DeleteSocialMedia [DELETE]

## Get Started
1. Clone this repository
2. Create and set up your own environment variables by using the **.env.example file**. Input your postgres database configuration and your jwt secret key as well as your server port
3. And now you're ready to go
4. Run this project `go run main.go`

## Package or Library Used
- Go Framework: **Gin**
- ORM: **GORM**
- Database: **PostgreSQL**
- Environment Variable: **Viper**
- API Documentation: **Swagger**