<template>
  <section class="container-fluid">
    <vue-recaptcha
      ref="recaptcha"
      :sitekey="siteKey"
      size="invisible"
      badge="right"
      @verify="onVerify"
      @expired="onExpired"
      :loadRecaptchaScript="true"
    >
    </vue-recaptcha>
    <section class="container">
      <div class="url__form" role="url">
        <div>
          <button class="btn btn-success btn-lg" disabled > Scan</button>
          <button class="btn btn-primary btn-lg" @click.prevent="monitorMode">Monitor</button>
        </div>
        <a href="/"><img id="header__title" src="../assets/images/title.png"/></a>
        <a href="/"><img id="header__logo" src="../assets/images/logo.png"/></a>
        <!-- <img id="header__text" src="../assets/images/text.png"/> -->
        <div class="col-md-12 text-center subtitle mt-3" >ANALYZE SUSPICIOUS URLS TO DETECT DEFACEMENT</div>
        <div v-if="processing" class="processing">
          <div class></div>
          <br>
          <div class="align-items-center justify-content-center">
            <div class="d-flex">
              <div class="lds-roller">
                <div></div>
                <div></div>
                <div></div>
                <div></div>
                <div></div>
                <div></div>
                <div></div>
                <div></div>
              </div>
              <h3 class="ml-2 mt-3">Processing</h3>

            </div>
          </div>
        </div>

        <div v-if="!processing" class="search__header">
          <input
            v-model="url"
            type="text"
            class="form-control input-lg"
            placeholder="Enter URL here"
            @keyup.enter="scanDeface"
          />
          <button class="btn btn-lg btn-primary btn-block" @click.prevent="scanDeface">
            <vue-octicon :icon="search">Search</vue-octicon>
          </button>
          <div class="alert alert-success" v-if="hasMessages ===2">
            <strong>Success!</strong> {{ message}}
          </div>
          <div class="alert alert-danger" v-if="hasMessages === 1">
            <strong>{{ message }}</strong>
          </div>
        </div>
      </div>

    </section>
  </section>
</template>


<script>
  import AnalysisAPI from '@/api'
  import VueRecaptcha from 'vue-recaptcha'
  import VueOcticon, {search} from 'octicons-vue'

  export default {
    data() {
      return {
        scanType: 'url',
        url: '',
        data: '',
        siteKey: process.env.VUE_APP_SITE_KEY || '',
        processing: false,
        search,
        hasMessages: 0,
        message: ''
      }
    },
    components: {
      VueOcticon,
      VueRecaptcha
    },
    computed: {
      api() {
        return AnalysisAPI
      },
      computedUrl() {
        if (this.url.indexOf('http://') === -1 && this.url.indexOf('https://') === -1) {
          return 'http://' + this.url.trim()
        }

        return this.url.trim()
      }
    },
    methods: {
      clear() {
        this.url = ''
        this.$refs.recaptcha.reset()
        this.hasMessages = 0
      },
      scanDeface() {
        if (!this.url) {
          return
        }

        this.scanType = 'url'

        this.processing = true

        this.$refs.recaptcha.execute()
      },
      onVerify(response) {
        if (this.scanType === 'url') {
          let requestData = JSON.stringify({url: this.computedUrl, recaptcha_token: response})
          this.api.scanUrl(requestData).then(data => {
            if (data.status && data.status !== 200 && data.status !== 201) {
              this.clear()

              this.processing = false

              return
            }

            this.$router.push({name: 'report-deface', params: {data: data}})

          }).catch((err) => {
            this.clear()
            this.hasMessages = 1
            if (!err.response) {
              return
            }

            this.message = err.response.data.message
            this.processing = false
          })

        }
      },
      onExpired() {
        this.$refs.recaptcha.reset()
      },
      monitorMode() {
        this.$router.push({name: 'monitor'})
      }
    }
  }
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style src="../assets/bootstrap/css/bootstrap.css"></style>
<style src="../assets/style.css" scoped>
</style>
