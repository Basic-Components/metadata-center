import Vue from 'vue'
import VueHighlightJS from 'vue-highlight.js'
// Highlight.js languages (Only required languages)
import json from 'highlight.js/lib/languages/json'

import 'highlight.js/styles/solarized-light.css'
Vue.use(VueHighlightJS, {
    languages: {
        json
    }
})