<template>
  <div>
    <Bar :data="chartData" :options="chartOptions" />
  </div>
</template>

<script setup>
import { Bar } from "vue-chartjs";

const props = defineProps(["options", "responses"]);

const chartData = computed(() => {
  if (props.options && props.responses) {
    const backgroundColors = [];
    const response = Object.keys(props.options)?.reduce(
      (responses, option) => {
        responses[option] = 0;

        if (props.options[option].isAnswer) {
          backgroundColors.push("#17B169");
        } else {
          backgroundColors.push("#fd5c63");
        }
        return responses;
      },
      { "Not Attempted": 0 }
    );
    backgroundColors.push("#FFCC00");

    const responses = props.responses?.reduce((acc, response) => {
      const { String: answer, Valid } = response.answers;

      if (!Valid || answer === "") {
        // Increment "Not Attempted" count if the answer is invalid or empty
        acc["Not Attempted"]++;
      } else {
        // Parse the answer (assuming the answer is in the format "[n]")
        const answerValue = parseInt(answer.replace(/\[|\]/g, ""), 10);
        if (acc[answerValue] !== undefined) {
          // Increment the count for the valid answer
          acc[answerValue]++;
        }
      }

      return acc;
    }, response);

    return {
      labels: Object.keys(responses),
      datasets: [
        {
          datalabels: {
            labels: {
              title: null,
            },
          },
          data: Object.values(responses),
          backgroundColor: backgroundColors,
          borderWidth: 1,
          borderRadius: 5,
          barThickness: 20,
        },
      ],
    };
  }
});
const chartOptions = ref({
  responsive: true,
  maintainAspectRatio: false,
  scales: {
    x: {
      grid: {
        display: false,
      },
      ticks: {
        callback: function (value, index) {
          const val = this.getLabelForValue(value);
          let xAxisName = props.options[val]?.value || "Not Attempted"
          if (xAxisName.length > 15) {
            xAxisName = xAxisName.slice(0,15) + "..."
          }
          return xAxisName;
        },
      },
    },
    y: {
      beginAtZero: true,
      ticks: {
        stepSize: 1,
      },
    },
  },
  plugins: {
    legend: {
      display: false,
    },
  },
});
</script>
