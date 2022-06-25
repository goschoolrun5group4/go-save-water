/* globals Chart:false, feather:false */

(function () {
  'use strict'

  feather.replace({ 'aria-hidden': 'true' })

  let ctx = document.getElementById('myChart');

  const userConsumption = document.getElementById('userConsumption').value;
  const nationalConsumption = document.getElementById('nationalConsumption').value;

  if (ctx != null) {
    const data = {
      labels: getMonths(),
      datasets: [{
        type: 'bar',
        label: 'Usage (Cu M)',
        data: JSON.parse(userConsumption),
        backgroundColor: 'rgb(54, 162, 235)',
        borderColor: 'rgb(54, 162, 235)',
        order: 2
      }, {
        type: 'line',
        tension: 0,
        label: 'National Average (Cu M)',
        data: JSON.parse(nationalConsumption),
        fill: false,
        borderColor: 'rgb(255, 99, 132)',
        order: 1
      }]
    };

    let myChart = new Chart(ctx, {
      type: 'bar',
      data: data,
      options: {
        layout: {
          padding: 30
        },
        scales: {
          yAxes: [{
            ticks: {
              beginAtZero: true
            }
          }]
        },
        legend: {
          display: true
        },
        title: {
          display: true,
          text: 'Consumption Trend'
        }
      }
    })

    function getMonths() {
      const month = ["January","February","March","April","May","June","July","August","September","October","November","December"];
      let dataMonth = [];

      for (let i = 5; i >= 0; i--) {
        let date = new Date()
        date.setMonth(date.getMonth() - i)
        dataMonth.push(month[date.getMonth()])
      }

      return dataMonth
    }
  }


})()
