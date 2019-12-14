<template>
  <section class="container-fluid">

    <section class="report">

      <div style="padding:0% 40%">
        <button class="btn btn-success btn-lg" @click.prevent="scanMode" > Scan</button>
        <button class="btn btn-primary btn-lg" @click.prevent="monitorMode" >Monitor</button>
      </div>
      <a href = '/'><img id="header__title" src="../assets/images/title.png"/></a>
      <a href = '/'><img id="header__logo" src="../assets/images/logo.png"/></a>
      <!-- <img id="header__text" src="../assets/images/text.png"/> -->
      <div class="col-md-12 text-center subtitle mt-3" >ANALYZE SUSPICIOUS URLS TO DETECT DEFACEMENT</div>


      <!-- REPORT DEFACEMENT -->

      <div v-if="result">

        <div class="search__header">
          <h1 >Report</h1>
        </div>

        <table width="100%" border="0">
          <tr>
            <td class="text-right" width="25%">
              <label>URL:</label>
            </td>
            <td width="80%">
              <input type="text" v-model="result.startURL" disabled/>
            </td>
          </tr>
        </table>

        <div v-if="!defaced" class="result__alert">
          <h3 class="text-success">
            <img src="../assets/images/safe.png" width="30px" height="30px"/> Benign
          </h3>
        </div>

        <div v-if="defaced" class="result__alert">
          <h3 class="text-danger">
            <img src="../assets/images/defaced.png" width="30px" height="30px"/> Defaced
          </h3>
        </div>


        <div class="result__chart">
          <h1 class>Convolutional Neural Network</h1>
          <div class="chart">
            <vue-apex-charts v-if="CNN && CNN.series && CNN.series.length" width="100%" type="donut"
                             :series="CNN.series" :options="chartOptions"></vue-apex-charts>
          </div>
        </div>


        <div class="result__screenshot">
          <h2>Screenshot</h2>
          <div style="outline:1px solid;margin:auto;width:900px;height:900px;">
            <img :src="base64Image" width="100%"/>
          </div>
        </div>
      </div>
    </section>
  </section>
</template>


<script>

  import VueApexCharts from 'vue-apexcharts'
  import AnalysisAPI from '@/api'
  import VueOcticon, {sync} from 'octicons-vue'
  import {setTimeout} from 'timers'

  export default {
    data() {
      return {
        siteKey: process.env.VUE_APP_SITE_KEY || '',
        reportId: '',
        processing: false,
        score: 0,
        CNN: {
          detected: false,
          series: [0.7, 0.3],
          score: 0
        },
        defaced: false,
        result: {
          startURL: '3x4mp13.c0m'
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
        if (!this.result || !this.result.data || !this.result.data.screenshot) {
          return
        }
        return `data:image/jpg;base64,${this.result.data.screenshot}`
      }
    },
    methods: {
      fetchData() {
        if (!this.$route.params.data) {
          window.location.href = '/'
          return
        }
        this.result = this.$route.params.data
        this.result.startURL = this.$route.params.data.data.url
        this.defaced = this.result.data.prediction === '1'
      },
      updateData() {
        if (!this.result ||
          !this.result.data ||
          !this.result.data.prediction
        ) {
          return
        }
        this.score = parseFloat(this.result.data.percentage)
        let CNNPercentage = Math.round(this.score * 100) / 100

        this.CNN = {
          detected: this.result.data.prediction === '1',
          series: [(1-CNNPercentage),
            CNNPercentage],
          score: this.score
        }

        this.score = this.CNN.score
        this.$nextTick(function () {
          window.scrollBy({
            top: 250,
            left: 0,
            behavior: 'smooth'
          })
        })
      }, removeFragment(url) {
        return url.split('#')[0]
      },
      removeQueryString(url) {
        return url.split('?')[0]
      },
      stripUrl(url) {
        let tmp = url
        tmp = this.removeFragment(tmp)
        tmp = this.removeQueryString(tmp)
        return tmp
      },
      monitorMode() {
        this.$router.push({name: 'monitor'})
      },
      scanMode() {
        this.$router.push({name: 'scan'})
      },
    },
    mounted() {
      this.fetchData()
      setTimeout(() => {
        this.updateData()
      }, 1000)
    }
  };
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style src="../assets/bootstrap/css/bootstrap.css"></style>
<style src="../assets/style.css" scoped></style>
