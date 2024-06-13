package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/charmbracelet/huh"
)

type Exercise struct {
    Name     string
    Sets     int
    Reps     []int
	Weights  []float64
    Rest     string
	SupersetWith string
}

type WorkoutSession struct {
    Date      time.Time
    Exercises []Exercise
	Name      string
}

type Improvement struct {
    Name           string
    WeightIncrease float64
    RepIncrease    int
}

func saveWorkouts(workouts []WorkoutSession, filename string) error {
    data, err := json.Marshal(workouts)
    if err != nil {
        return err
    }
    return ioutil.WriteFile(filename, data, 0644)
}

func loadWorkouts(filepath string) ([]WorkoutSession, error) {
    if _, err := os.Stat(filepath); os.IsNotExist(err) {
        fmt.Printf("File does not exist: %s\n", filepath)
        return nil, err
    }

    data, err := ioutil.ReadFile(filepath)
    if err != nil {
        return nil, err
    }

    var loadedWorkouts []WorkoutSession
    if err := json.Unmarshal(data, &loadedWorkouts); err != nil {
        return nil, err
    }

    return loadedWorkouts, nil
}

func inputExerciseDetails(exercise *Exercise, completedExercises map[string]bool, allExercises []Exercise) error {
	if _, done := completedExercises[exercise.Name]; done {
		return nil
	}

	// Input reps for each set
	exercise.Reps = make([]int, exercise.Sets)
    for j := 0; j < exercise.Sets; j++ {
        for {
            fmt.Print("\033[H\033[2J") // Clear the terminal
            prompt := fmt.Sprintf("Input details for %s:\nSet %d - Enter reps:", exercise.Name, j+1)
            var reps string
            err := huh.NewInput().Title(prompt).Value(&reps).Run()
            if err != nil {
                return err
            }
            exercise.Reps[j], err = strconv.Atoi(reps)
            if err != nil {
                continue
            }
            break
        }
    }

	// Input weights for each set
    exercise.Weights = make([]float64, exercise.Sets)
    for j := 0; j < exercise.Sets; j++ {
        for {
            fmt.Print("\033[H\033[2J") // Clear the terminal
            prompt := fmt.Sprintf("Input details for %s:\nSet %d - Enter weight:", exercise.Name, j+1)
            var weight string
            err := huh.NewInput().Title(prompt).Value(&weight).Run()
            if err != nil {
                return err
            }
            exercise.Weights[j], err = strconv.ParseFloat(weight, 64)
            if err != nil {
                continue
            }
            break
        }
    }

	// Handle supersets
	if exercise.SupersetWith != "" {
		fmt.Printf("Proceed to superset exercise: %s\n", exercise.SupersetWith)
	}

	// Specific logic to skip the counterpart exercise if needed
	if exercise.Name == "Benkpress med stong" {
		completedExercises["Benkpress med manuala"] = true
	} else if exercise.Name == "Benkpress med manuala" {
		completedExercises["Benkpress med stong"] = true
	}

	return nil
}

func logWorkout(workouts, loggedWorkouts []WorkoutSession) {
    fmt.Print("\033[H\033[2J")
    workoutOptions := make([]huh.Option[string], len(workouts)+1)  // Include space for 'Back to Main Menu'
    for i, workout := range workouts {
        workoutOptions[i] = huh.NewOption(workout.Name, workout.Name)
    }
    // Add 'Back to Main Menu' option at the end
    workoutOptions[len(workouts)] = huh.NewOption("Return to main menu", "Return to main menu")

    var selectedWorkoutName string
    err := huh.NewSelect[string]().Title("Choose your workout or go back:").Options(workoutOptions...).Value(&selectedWorkoutName).Run()
    if err != nil {
        fmt.Println("Error selecting workout:", err)
        return
    }

    if selectedWorkoutName == "Return to main menu" {
        mainMenu(workouts, loggedWorkouts) // Return to main menu if selected
        return
    }

    var selectedWorkout *WorkoutSession
    for _, workout := range workouts {
        if workout.Name == selectedWorkoutName {
            selectedWorkout = &workout
            break
        }
    }

    if selectedWorkout == nil {
        fmt.Println("Selected workout not found")
        return
    }

    selectedWorkout.Date = time.Now()
    completedExercises := make(map[string]bool)

    // Pre-determine if both exercises are in the workout
    hasStong, hasManuala := false, false
    for _, ex := range selectedWorkout.Exercises {
        if ex.Name == "Benkpress med stong" {
            hasStong = true
        } else if ex.Name == "Benkpress med manuala" {
            hasManuala = true
        }
    }

    // Prompt the user to choose one if both are present
    if hasStong && hasManuala {
        var choice string
        options := []huh.Option[string]{
            huh.NewOption("Benkpress med stong", "Benkpress med stong"),
            huh.NewOption("Benkpress med manuala", "Benkpress med manuala"),
        }
        err = huh.NewSelect[string]().Title("Choose type of Benkpress:").Options(options...).Value(&choice).Run()
        if err != nil {
            fmt.Println("Error choosing between exercises:", err)
            return
        }
        if choice == "Benkpress med stong" {
            completedExercises["Benkpress med manuala"] = true
        } else {
            completedExercises["Benkpress med stong"] = true
        }
    }

    // Proceed with workout logging
    for _, ex := range selectedWorkout.Exercises {
        if !completedExercises[ex.Name] {
            err := inputExerciseDetails(&ex, completedExercises, selectedWorkout.Exercises)
            if err != nil {
                fmt.Println("Error inputting exercise details:", err)
                return
            }
        }
    }

    // Save the workout data
    err = saveWorkouts(loggedWorkouts, "log.json")
    if err != nil {
        fmt.Println("Failed to save workouts:", err)
    }

    mainMenu(workouts, loggedWorkouts)
}

func showProgressionInteractive(workouts, loggedWorkouts []WorkoutSession) {
    now := time.Now()
    oneMonthAgo := now.AddDate(0, -1, 0)
    oneWeekAgo := now.AddDate(0, 0, -7)

    totalWorkouts, workoutsLastMonth, workoutsLastWeek := 0, 0, 0
    improvements := []Improvement{}

    for _, workout := range loggedWorkouts {
        if workout.Date.After(oneMonthAgo) {
            workoutsLastMonth++
        }
        if workout.Date.After(oneWeekAgo) {
            workoutsLastWeek++
        }

        totalWorkouts++
        imp := calculateImprovements(workout, loggedWorkouts)
        improvements = append(improvements, imp...)
    }

    sort.Slice(improvements, func(i, j int) bool {
        return improvements[i].WeightIncrease > improvements[j].WeightIncrease
    })

    totalSummary := fmt.Sprintf("Total workouts logged: %d\nWorkouts in the last month: %d\nWorkouts in the last week: %d\n", totalWorkouts, workoutsLastMonth, workoutsLastWeek)

    fmt.Println("Progression summary:")
    fmt.Println(totalSummary)
    fmt.Println("\nTop Improvements:")

    improvementDetails := generateImprovementFields(improvements)
    for _, detail := range improvementDetails {
        fmt.Println(detail)
    }

    // Post-view action selection
    var action string
    err := huh.NewSelect[string]().
        Title("What would you like to do next?").
        Options(
            huh.NewOption("Return to main menu", "menu"),
            huh.NewOption("Quit", "quit"),
        ).
        Value(&action).
        Run()

    if err != nil {
        fmt.Println("Error:", err)
        return
    }

    if action == "menu" {
        mainMenu(workouts, loggedWorkouts)
    } else if action == "quit" {
        fmt.Println("Exiting program.")
        os.Exit(0)
    }
}

func generateImprovementFields(improvements []Improvement) []huh.Field {
    fields := []huh.Field{}
    for i, imp := range improvements {
        if i >= 3 {  // Limiting to top 3 for simplicity
            break
        }
        detail := fmt.Sprintf("%d. %s: +%0.2fkg", i+1, imp.Name, imp.WeightIncrease)
        fmt.Println(detail)
    }
    return fields
}

func calculateImprovements(workout WorkoutSession, workouts []WorkoutSession) []Improvement {
    var improvements []Improvement
    for _, exercise := range workout.Exercises {
        // Find the first and last instances of each exercise within the timeframe
        firstInstance := findFirstInstance(exercise.Name, workouts, time.Now().AddDate(0, -1, 0))
        if firstInstance == nil || len(firstInstance.Weights) == 0 || len(exercise.Weights) == 0 {
            continue // Skip if no first instance or weights are recorded incorrectly
        }

        // Assuming last instance is the current one since we're iterating from latest to earliest
        lastWeight := exercise.Weights[len(exercise.Weights)-1]
        firstWeight := firstInstance.Weights[0]
        lastReps := exercise.Reps[len(exercise.Reps)-1]
        firstReps := firstInstance.Reps[0]

        weightIncrease := lastWeight - firstWeight
        repIncrease := lastReps - firstReps

        if weightIncrease > 0 || repIncrease > 0 {
            improvements = append(improvements, Improvement{
                Name:           exercise.Name,
                WeightIncrease: weightIncrease,
                RepIncrease:    repIncrease,
            })
        }
    }
    return improvements
}

func findFirstInstance(name string, workouts []WorkoutSession, since time.Time) *Exercise {
    earliestWorkoutDate := time.Now() // Assuming we use the current time as the initial comparison point
    var earliestExercise *Exercise

    for _, workout := range workouts {
        if workout.Date.Before(since) || workout.Date.After(earliestWorkoutDate) {
            continue
        }
        for _, exercise := range workout.Exercises {
            if exercise.Name == name {
                if earliestExercise == nil || workout.Date.Before(earliestWorkoutDate) {
                    earliestWorkoutDate = workout.Date
                    // Make a copy of the exercise because we are going to return a reference to it
                    copy := exercise
                    earliestExercise = &copy
                }
            }
        }
    }
    return earliestExercise
}

func mainMenu(workouts, loggedWorkouts []WorkoutSession) {
    fmt.Print("\033[H\033[2J")

    var choice string
    err := huh.NewSelect[string]().
        Title("Welcome to Helsi! What would you like to do?").
        Options(
            huh.NewOption("Log new workout", "Log new workout"),
            huh.NewOption("Show progression", "Show progression"),
            huh.NewOption("Quit", "Quit"),
        ).
        Value(&choice).
        Run()

    if err != nil {
        fmt.Println("Error:", err)
        return
    }

    // Execute the chosen option
    switch choice {
    case "Log new workout":
        logWorkout(workouts, loggedWorkouts) // Pass both slices
    case "Show progression":
        showProgressionInteractive(workouts, loggedWorkouts) // Use only logged workouts
    case "Quit":  // Handling the Quit option
        fmt.Println("Exiting program.")
        os.Exit(0)
    }
}

func main() {
    fmt.Print("\033[H\033[2J")

    configPath := flag.String("config", "workouts.json", "path to the configuration file containing workouts")
    flag.Parse()

    workouts, err := loadWorkouts(*configPath)
    if err != nil {
        fmt.Printf("Failed to load workouts from file: %v\n", err)
        os.Exit(1)
    }
    loggedWorkouts, err := loadWorkouts("log.json") // Load logged workouts separately
    if err != nil {
        loggedWorkouts = []WorkoutSession{} // Initialize if not available
    }
    mainMenu(workouts, loggedWorkouts)
}
