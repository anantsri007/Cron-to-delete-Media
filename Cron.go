package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func main() {
	dirpath := `C:\Users\hp\Documents\Emptyfolder1\folder6`

	// Ask the user for the number of days to subtract
	var daysStr string
	fmt.Println("Enter the number of days to subtract:")
	fmt.Scanln(&daysStr)

	// Convert input to integer
	days, err := strconv.Atoi(daysStr)
	if err != nil || days < 0 {
		fmt.Println("Invalid number of days. Please enter a positive integer.")
		return
	}

	// Get toy's date and calculate the target date
	today := time.Now()
	targetDate := today.AddDate(0, 0, -days)
	fmt.Printf("Target date for deluuetion: %s\n", targetDate.Format("02 Jan 2006"))

	// Iterate over the directory to delete files created on the target date
	err = iterateAndDeleteFilesAndEmptyDirs(dirpath, targetDate)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Files deleted and empty directories cleaned up.")
	}
}

// iterateAndDeleteFilesAndEmptyDirs walks through the directory, deletes files created on the target date, and deletes empty folders.
func iterateAndDeleteFilesAndEmptyDirs(dirpath string, targetDate time.Time) error {
	return filepath.Walk(dirpath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// If it's a file, check its creation/modification date
		if !info.IsDir() {
			// Get the file modification time
			fileModTime := info.ModTime()

			// Check if the file was modified/created on the target date
			if isSameDay(fileModTime, targetDate) {
				// Delete the file
				err := os.Remove(path)
				if err != nil {
					return fmt.Errorf("could not delete file %s: %w", path, err)
				}
				fmt.Printf("Deleted file: %s\n", path)
			}
		}

		// If it's a directory (and not the root directory), check if it's empty
		if info.IsDir() && path != dirpath {
			isEmpty, err := isDirEmpty(path)
			if err != nil {
				return err
			}

			// If the directory is empty, delete it
			if isEmpty {
				err := os.Remove(path)
				if err != nil {
					return fmt.Errorf("could not delete directory %s: %w", path, err)
				}
				fmt.Printf("Deleted empty directory: %s\n", path)
			}
		}
		return nil
	})
}

// isDirEmpty checks whether a directory is empty.
func isDirEmpty(path string) (bool, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return false, err
	}

	return len(files) == 0, nil
}

// isSameDay checks if two time values represent the same calendar day.
func isSameDay(t1, t2 time.Time) bool {
	year1, month1, day1 := t1.Date()
	year2, month2, day2 := t2.Date()

	return year1 == year2 && month1 == month2 && day1 == day2
}
