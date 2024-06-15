<template>
  <div class="container">
    <h1>Log workout</h1>
    <select v-model="selectedWorkout" @change="fetchExercises">
      <option v-for="workout in workouts" :key="workout.Name" :value="workout.Name">
        {{ workout.Name }}
      </option>
    </select>

    <input type="date" v-model="workoutDate" placeholder="Workout Date">

    <div class="exercises" v-for="exercise in filteredExercises" :key="exercise.Name">
      <h2>{{ exercise.Name }}</h2>
      <div class="set-details" v-for="setIndex in exercise.Sets" :key="setIndex">
        Set {{ setIndex }}:
        <input type="number" v-model="exercise.Reps[setIndex-1]" @change="handleInput(exercise.Name)" placeholder="Reps"/>
        <input type="number" v-model="exercise.Weights[setIndex-1]" @change="handleInput(exercise.Name)" placeholder="Weight (kg)"/>
      </div>
    </div>

    <button @click="saveWorkout">Save Workout</button>
  </div>
</template>

<script>
import axios from 'axios'

export default {
  data() {
    return {
      workouts: [],
      exercises: [],
      selectedWorkout: null,
      hideExercises: {},
      workoutDate: ''
    };
  },
  created() {
    this.loadWorkouts();
  },
  computed: {
    filteredExercises() {
      return this.exercises.filter(exercise => !this.hideExercises[exercise.Name]);
    }
  },
  methods: {
    loadWorkouts() {
    axios.get('/workouts').then(response => {
      this.workouts = response.data;
    }).catch(error => {
      console.error("Error loading workouts:", error);
      alert('Failed to load workouts: ' + error.message);
    });
  },
  fetchExercises() {
    const workout = this.workouts.find(w => w.Name === this.selectedWorkout);
    this.exercises = workout ? workout.Exercises : [];
    this.disabledExercises = {};
  },
  handleInput(exerciseName) {
    const currentExercise = this.exercises.find(ex => ex.Name === exerciseName);
    const weightsNotEmpty = currentExercise.Weights.some(weight => weight > 0);

    if (exerciseName === "Benkpress med stong") {
      this.hideExercises["Benkpress med manuala"] = weightsNotEmpty;
      this.hideExercises["Benkpress med stong"] = !weightsNotEmpty;
    } else if (exerciseName === "Benkpress med manuala") {
      this.hideExercises["Benkpress med stong"] = weightsNotEmpty;
      this.hideExercises["Benkpress med manuala"] = !weightsNotEmpty;
    }

    if (!weightsNotEmpty) {
      this.hideExercises["Benkpress med stong"] = false;
      this.hideExercises["Benkpress med manuala"] = false;
    }
  },
  saveWorkout() {
      if (!this.workoutDate) {
        alert('Please select a workout date.');
        return;
      }

      const workoutData = {
        Name: this.selectedWorkout,
        Date: this.workoutDate,
        Exercises: this.exercises
      };

      console.log("Saving workout session:", workoutData);
      axios.post('/log', workoutData)
        .then(() => {
            alert('Workout saved!');
        })
        .catch(error => {
            console.error("Error saving workout:", error);
            console.error("Error details:", error.response.data);
        });
    }
}

};
</script>

<style lang="scss">
$background-color: #000;
$text-color: #fff;
$button-color: #673ab7;
$input-background: #9c27b0;
$input-text-color: #fff;

.container {
  background-color: $background-color;
  color: $text-color;
  padding: 20px;
  border-radius: 8px;
  width: 90%;
  max-width: 600px;
  margin: 20px auto;
  box-shadow: 0 4px 6px rgba(0,0,0,0.1);
}

select, input[type="number"] {
  width: 100%;
  padding: 8px;
  margin: 10px 0;
  background-color: $input-background;
  color: $input-text-color;
  border: 1px solid darken($input-background, 15%);
  border-radius: 4px;
  outline: none;
  box-shadow: 0 2px 4px rgba(0,0,0,0.2);

  &::placeholder {
    color: lighten($input-text-color, 20%);
  }
}

button {
  width: 100%;
  padding: 10px;
  background-color: $button-color;
  color: $text-color;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.3s ease;

  &:hover {
    background-color: lighten($button-color, 10%);
  }
}

.exercises {
  margin-top: 20px;
}

.exercise-sets {
  display: flex;
  justify-content: space-between;
  input {
    flex: 1;
    &:not(:last-child) {
      margin-right: 10px;
    }
  }
}
</style>
