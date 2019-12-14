import Vue from 'vue'
import Router from 'vue-router'
import Scan from '@/components/Scan'
import ReportDeface from '@/components/ReportDeface'
import Monitor from '@/components/Monitor'

Vue.use(Router)

export default new Router({
    routes: [{
            path: '/',
            name: 'scan',
            component: Scan
        },
        {
            path: '/report',
            name: 'report-deface',
            component: ReportDeface
        },
        {
          path: '/monitor',
          name: 'monitor',
          component: Monitor
        }
    ]
})
