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
        <a href="/"><img id="header__title" src="../assets/images/title.png"/></a>
        <a href="/"><img id="header__logo" src="../assets/images/logo.png"/></a>
        <!-- <img id="header__text" src="../assets/images/text.png"/> -->
        <div class="col-md-12 text-center subtitle mt-3" >MONITORING URLS TO DETECT DEFACEMENT</div>
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
            @keyup.enter="monitorDeface"
          />
          <input
            v-model="email"
            type="text"
            class="form-control input-lg"
            placeholder="Enter email to send alert here"
            @keyup.enter="monitorDeface"
          />
          <div class="form-group">
            <label>Every</label>
            <input
            v-model="interval"
            type="number"
            class="input-lg"
            placeholder="Enter interval to monitor here"
            @keyup.enter="monitorDeface"
          />
            <label>minute(s)</label>
          </div>
          <button class="btn btn-lg btn-primary btn-block" @click.prevent="monitorDeface">
            <vue-octicon :icon="plus">Add</vue-octicon>
          </button>
          <div class="alert alert-success" v-if="hasMessages ===2">
            <strong>Success!</strong> {{ message}}
          </div>
          <div class="alert alert-danger" v-if="hasMessages === 1">
            <h4>{{ message }}</h4>
          </div>
        </div>
      </div>
    </section>
  </section>
</template>

<script>
    import AnalysisAPI from '@/api'
    import VueRecaptcha from 'vue-recaptcha'
    import VueOcticon, { plus } from 'octicons-vue'
    export default {
        name: "Monitor",
        data() {
          return {
            scanType: 'url',
            url: '',
            data: '',
            siteKey: process.env.VUE_APP_SITE_KEY || '',
            processing: false,
            plus,
            hasMessages: 0,
            message: '',
            email: '',
            interval: 5,
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
            this.message = ''
            this.email = ''
            this.interval = 5
            this.$refs.recaptcha.reset()
            this.hasMessages = 0
          },
          monitorDeface() {
            if (!this.url || !this.email) {
              return
            }

            this.scanType = 'url'

            this.processing = true

            this.$refs.recaptcha.execute()
          },
          onVerify(response) {
            if (this.scanType === 'url') {
              let requestData = JSON.stringify({url: this.computedUrl, recaptcha_token: response, email: this.email, interval: parseInt(this.interval)})
              this.api.monitor(requestData).then(data => {
                this.clear()
                this.hasMessages = 2
                this.processing = false
                this.message = data.message
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
          }
        }

    }
</script>

<style src="../assets/bootstrap/css/bootstrap.css"></style>
<style src="../assets/style.css" scoped>
</style>

