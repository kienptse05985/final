<template>
  <section class="container-fluid">
    <section class="report">
      <div class="bg-rotate"></div>

      <img id="header__title" src="../assets/images/title.png" />
      <img id="header__logo" src="../assets/images/logo.png" />
      <img id="header__text" src="../assets/images/text.png" />


      <!-- REPORT DEFACEMENT -->

      <div v-if="result">

      <div class="search__header">
        <h2 class>Report</h2>
      </div>

      <table width="100%" border="0">
        <tr>
          <td class="text-right" width="35%">
            <label>URL:</label>
          </td>
          <td width="80%">
            <input type="text" v-model="result.startURL" disabled />
          </td>
        </tr>
          <tr v-if="stripUrl(result.destination) !== stripUrl(result.startURL)">
          <td class="text-right" width="20%">
            <label>Final Redirected URL:</label>
          </td>
          <td width="80%">
            <input type="text" v-model="result.destination" disabled />
          </td>
        </tr>
      </table>

      <div v-if="score === 0" class="result__alert">
        <h3 class="text-success">
          <img src="../assets/images/safe.png" width="30px" height="30px" /> Safe
        </h3>
      </div>

      <div v-if="score === 0" class="result__alert">
        <h3 class="text-danger">
          <img src="../assets/images/defaced.png" width="30px" height="30px" /> Danger
        </h3>
      </div>

      <div class="result__chart">
        <div class="chart">
          <h6 class>Logistic Regression</h6>
          <vue-apex-charts v-if="logisticRegression && logisticRegression.series && logisticRegression.series.length" width="400px" type="donut" :series="logisticRegression.series" :options="chartOptions"></vue-apex-charts>
        </div>

        <div class="chart">
          <h6 class>Random Forest</h6>
          <vue-apex-charts v-if="randomForest && randomForest.series && randomForest.series.length" width="400px" type="donut" :series="randomForest.series" :options="chartOptions"></vue-apex-charts>
        </div>
        
      </div>
      <div v-if="!processing" class="result__screenshot">
        <h2>Screenshot</h2>

        <div class>
          <img :src="base64Image" width="100%" />
        </div>
        <div style="outline:1px solid;margin:auto;width:900px;height:900px;"></div>
      </div>
      </div>
    </section>
  </section>
</template>


<script>

import VueApexCharts from 'vue-apexcharts'
import AnalysisAPI from '@/api'
import VueOcticon, {sync} from 'octicons-vue'
import { setTimeout } from 'timers'

export default {
  data() {
    return {
      siteKey: process.env.VUE_APP_SITE_KEY || '',
      reportId: '',
      processing: false,
      result: null,
      score: 0,
      logisticRegression: {
        detected: false,
        series: [0.7, 0.3],
        score: 0
      },
      result: {
        startURL: '3x4mp13.c0m',
        destination: 'http://3x4mp13.c0m'
      },
      randomForest: {
        detected: false,
        series: [0.25, 0.75],
        score: 0
      },
      chartOptions: {
        labels: ['Defacement', 'Clean'],
        colors: ['#f54560', '#21d77d'],
        legend: {
          position: 'bottom',
          fontSize: '16px'
        }
      },
      timeOut: 0,
      error: {
        status: false,
        message: ''
      },
      sync
    }
  },
  components: {
    VueOcticon,
    VueApexCharts
  },
  computed: {
    api() {
      return AnalysisAPI
    },
    base64Image() {
      if (!this.result || !this.result.screenshot) {
        return
      }

      return `data:image/jpg;base64,${this.result.screenshot}`
    }
  },
  methods: {
    fetchData () {
      if (!this.$route.params.id) {
        return
      }

      this.reportId = this.$route.params.id
      this.processing = true
    },
    getReportDeface () {
      if (!this.reportId) {
        return
      }

      this.timeOut += 5000
      if (this.timeOut === 30000) {
        this.processing = false
        return
      }

      let requestData = JSON.stringify({ report_id: this.reportId })

      this.api.getUrlReport(requestData).then(data => {
        if (!(data && data.responseCode)) {

          this.error.status = true

          this.processing = false

          return
        }
        if (data.responseCode !== 200) {
          if (data.responseCode !== 202) {

            this.error.status = true
            this.error.message = 'Cannot get report'

            this.processing = false

            return
          }

          setTimeout(this.getReportDeface, 5000)

          return
        }

        if (data.responseCode === 200) {

          this.result = data

          this.updateData()
          this.processing = false
        }
      }).catch(error => {
        this.error.status = true
        this.processing = false
      })
    },
    updateData () {
      if (!this.result ||
        !this.result.models ||
        !this.result.models.logisticRegression ||
        !this.result.models.randomForest
      ) {
        return
      }

      let logisticRegressionPercentag = Math.round(this.result.models.logisticRegression.phishPercentage * 100) / 100
      let randomForestPercentage = Math.round(this.result.models.randomForest.phishPercentage * 100) / 100

      this.logisticRegression = {
        detected: this.result.models.logisticRegression.detected,
        series: [logisticRegressionPercentag,
          1 - logisticRegressionPercentag],
        score: this.result.models.logisticRegression.score
      }

      this.randomForest = {
        detected: this.result.models.randomForest.detected,
        series: [randomForestPercentage,
          1 - randomForestPercentage],
        score: this.result.models.randomForest.score
      }

      this.score = this.logisticRegression.score + this.randomForest.score
      this.$nextTick(function () {
        window.scrollBy({
          top: 300,
          left: 0,
          behavior: 'smooth'
        })
      })
    },removeFragment (url) {
      return url.split('#')[0]
    },
    removeQueryString (url) {
      return url.split('?')[0]
    },
    stripUrl (url) {
      let tmp = url
      tmp = this.removeFragment(tmp)
      tmp = this.removeQueryString(tmp)
      return tmp
    }
  },
  mounted() {
    this.fetchData()
    setTimeout(() => {
      this.getReportDeface()
    }, 1000)
  }
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style src="../assets/bootstrap/css/bootstrap.css"></style>
<style src="../assets/style.css" scoped></style>
