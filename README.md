# Helsi

CLI application to log workouts and view progression.

## How it works

Load predefined workouts from a `workouts.json` file in the following format:

```
[
    {
        "Name": "Workout 1",
        "Exercises": [
            {
                "Name": "Exercise 1",
                "Sets": 3,
                "Reps": [12, 10, 8],
                "Weights": [25, 27.5, 30],
                "Rest": "1.5 min",
                "SupersetWith": ""
            }
        ]
    }
]
```

Saved workouts will be logged to a separate `log.json` file

The app provides a CLI interface using [Huh](https://github.com/charmbracelet/huh) for prompts.