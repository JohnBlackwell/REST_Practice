## Password Validator

###
Checks to ensure that password adheres to the following standards:

  NIST recently updated their Digital Identity Guidelines in June 2017. The new guidelines specify general rules for handling the security   of user supplied passwords. Previously passwords were suggested to have certain composition rules (special characters, numbers, etc),       hints and expiration times. Those have gone out the window and the new suggestions are as follows:

  Passwords MUST

  Have an 8 character minimum
  AT LEAST 64 character maximum
  Allow all ASCII characters and spaces (unicode optional)
  Not be a common password
  
  ## Validator REST API
  
  ### Return Values
  
  All endpoints should return one of the following codes.
    
    200 Successful
    201 Created
    400 StatusBadRequest 
    403 Forbidden
    500 InternalServerError
  
  ### Process Endpoint
    
    /process   
        POST    //Takes either multi-part/form data ex. a text file containing common passwords
                 // or takes application/json and returns 201 if password passes, 401 if it doesn't
                 // returns 400 BadRequest if one of these two forms not sent in
                 
  ### Install
    git clone REST_Practice repository
  ### Usage
    ./buildDocker.sh inside cloned REST_Practice directory
   
    The AWS RDS Postgres DB is already seeded with 800k common passwords. To add more:
      curl -X POST -d @/path/to/password/file/   http://localhost:3000/process
     
     To validate a potential password:
     
     curl -d '{"password":"[password_to_test]"}' -H "Content-Type: application/json " -X POST http://localhost:3000/process
                 
     docker stop idea_container : ends server
     docker rm idea_container : removes idea_container. Do this when done with testing session. Rerun build script to rebuild container
  
  
