<template>
  <div class="hello">
    <h2>游戏便当steam代理</h2>
    <ul>
        <el-button size="large" :loading="loading" type="primary" @click="trigger">{{status ? '停止代理' : '启动代理'}}</el-button>
    </ul>
  </div>
</template>

<script>
import 'element-ui/packages/theme-chalk/lib/index.css';
import Vue from 'vue';
import {Button, Message} from 'element-ui';
import axios from 'axios';

Vue.use(Button);

Vue.prototype.$msg = Message;

export default {
  name: 'Main',
  data () {
    return {
        status: true,
        loading: false,
    }
  },
  mounted (){
      this.start();
  },
  methods: {
    start () {
        this.getstatus();
        return;
        setTimeout(() => {
            this.start();
        }, 1000);
    },
    getstatus () {
        this.loading = true;
        axios.get('/api/service/getstatus').then(resp => {
            resp = resp.data;
            this.loading = false;
            if(!resp.status){
                this.$msg.error(resp.message);
                return;
            }
            this.status = resp.data;
        }).catch(err => {
            console.log(err);
            this,$msg.error(err.message);
            this.loading = false;
        });
       
    },
    trigger (){
        this.loading = true;
        axios.get('/api/service/status/' + !this.status).then(resp => {
            this.loading = false;
            this.status = !this.status;
            return;
            resp = resp.data;
            if(!resp.status){
                this.$msg.error(resp.message); 
                return;
            }
        }).catch(err => {
            console.log(err); 
            this.$msg.error(err.message);
            this.loading = false;
       }); 
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
h1, h2 {
  font-weight: normal;
}
ul {
  list-style-type: none;
  padding: 0;
}
li {
  display: inline-block;
  margin: 0 10px;
}
a {
  color: #42b983;
}
</style>
