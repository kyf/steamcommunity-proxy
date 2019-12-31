<template>
  <div class="page-steam">
    <div class="hero">
      <div
          class="bg"
          :style="{
            backgroundImage: `url(${images.bg})`,
          }"></div>
      <div class="content">
        <h1>
          <img
              class="logo"
              alt="游戏便当"
              :src="images.logo">
          <span>Steam<br>社区加速</span>
        </h1>
        <button
            @click="trigger">
          <span
              v-if="!status">开始<br>加速</span>
          <span
              v-else>停止<br>加速</span>
        </button>
        <ul class="links">
          <li>
            <a @click="addCA">加载证书</a>
          </li>
          <li>
            <a href="https://youxibd.com/mycenter/setting">绑定Steam</a>
          </li>
        </ul>
      </div>
    </div>
    <div class="intro">
      <h2>操作说明</h2>
      <p>1.如果遇到Steam绑定页面无法打开的情况，请点击 <a href="">加载证书</a> 重试</p>
      <p>2.不需要Steam社区加速时，请点击 <a href="">停止加速</a> 如果页面不慎关闭，也可去右下角系统托盘区，右键点击加速图标进行操作，如下图所示：</p>
      <img
          :src="images.steps">
    </div>
  </div>
</template>

<script>
import bg from '../assets/bg.jpg';
import logo from '../assets/logo.png';
import steps from '../assets/steps.jpg';
import axios from 'axios';

export default {
  name: 'Steam',
  data() {
    return {
      images: {
        bg,
        logo,
        steps,
      },
      status: true,
    };
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
    addCA (){
        axios.get('/api/service/addca').then(resp => {
            resp = resp.data;
            this.loading = false;
            if(!resp.status){
                alert(resp.message);
                return;
            }
        }).catch(err => {
            console.log(err);
            alert(err.message);
            this.loading = false;
        });
       
    },
    getstatus () {
        this.loading = true;
        axios.get('/api/service/getstatus').then(resp => {
            resp = resp.data;
            this.loading = false;
            if(!resp.status){
                alert(resp.message);
                return;
            }
            this.status = resp.data;
        }).catch(err => {
            console.log(err);
            alert(err.message);
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
                alert(resp.message); 
                return;
            }
        }).catch(err => {
            console.log(err); 
            alert(err.message);
            this.loading = false;
       }); 
    }
   
  },
};
</script>

<style
    lang="scss"
    scoped
    src="./index.scss"></style>
