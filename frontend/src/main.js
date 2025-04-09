import Vue from "vue";
import App from "./App.vue";
import router from "./router";
import store from "./store";
import axios from "./service/api";
import ElementUI from "element-ui";
import "element-ui/lib/theme-chalk/index.css";
import mavonEditor from "mavon-editor";
import "mavon-editor/dist/css/index.css";
import "github-markdown-css/github-markdown.css";

Vue.use(ElementUI);
Vue.use(mavonEditor);

Vue.prototype.Free = window.Free;
Vue.prototype.$axios = axios;
Vue.config.productionTip = false;

router.beforeEach((to, from, next) => {
  console.log(to);
  console.log(from);
  if (to.meta.requireAuth) {
    if (localStorage.getItem("loginResult")) {
      next();
    } else {
      if (to.path === "/login") {
        next();
      } else {
        next({
          path: "/login",
        });
      }
    }
  } else {
    next();
  }

  if (to.fullPath == "/login") {
    if (localStorage.getItem("loginResult")) {
      next({
        path: from.fullPath,
      });
    } else {
      next();
    }
  }
});

router.afterEach(() => {
  document.title = "Gophers";
});

new Vue({
  router,
  store,
  render: (h) => h(App),
}).$mount("#app");
