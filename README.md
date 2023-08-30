# ThriftyLetterboxdCLI

This project was created for educational purposes and is not affiliated with Letterboxd. The project utilizes web scraping techniques due to the Letterboxd API being in private beta. It provides the ability to select random films from a watchlist and to view user's watched films stats.

## Disclaimer

This project is developed independently and is not endorsed or supported by Letterboxd. The purpose of this project is solely educational.

## Usage

Please note that this project is not intended for commercial use.

## How to Use

To use the ThriftyLetterboxd CLI, follow these steps:

1. Install Go (if not already installed): https://golang.org/dl/

2. Clone or download this repository.

3. Open a terminal and navigate to the project directory.

4. Build the executable:

   ```powershell
   go build .

   ```

5. Running the executeable.

   - On Windows:

     ```sh
     .\ThriftyLetterboxdCLI.exe stats <LetterboxdUsername>
     .\ThriftyLetterboxdCLI.exe random <LetterboxdUsername> [genre]
     ```

   - On Linux/macOS:
     ```sh
     ./ThriftyLetterboxdCLI stats <LetterboxdUsername>
     ./ThriftyLetterboxdCLI random <LetterboxdUsername> [genre]
     ```

   Replace `<LetterboxdUsername>` with the Letterboxd username you want to gather stats for, and `[genre]` with an optional genre keyword.
