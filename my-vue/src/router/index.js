import Vue from 'vue'
import Router from 'vue-router'
import Scan from '@/components/Scan'
import ReportDeface from '@/components/ReportDeface'

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
        }
    ]
})