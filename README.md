# Go Search Your File in Web

This is a simple web application built in Go that allows users to search for files in a specified directory using regular expressions.


## Features

- User-friendly HTML form for input.
- Recursive file search based on regex pattern.
- Displays search results on a separate results page.


## Setup and Running

1. **Clone the repository:**
    ```sh
   git clone https://github.com/Vkanhan/go-search-your-files.git
   cd go-search-your-files
   ```

2. **Run the application:**
    ```go
    go run main.go
    ```

3. **Access the application:**
    Open your web browser and go to http://localhost:8080 to access the search form.

## Usage 

    
   1. Open the application in your web browser at http://localhost:8080.

   2. Enter a regex pattern and a valid directory path in the search form.

   3. Submit the form to view the search results.


## Note
Ensure that the directory path provided in the search form is an absolute path and has the necessary read permissions.

## License
This project is licensed under the MIT [LICENSE](LICENSE). See the LICENSE file for details.


