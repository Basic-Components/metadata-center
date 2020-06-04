<template>
  <div class="service-new">
    <el-row type="flex" justify="center">
      <h2>创建服务</h2>
    </el-row>
    <el-row type="flex" justify="center" v-if="alertVisible">
      <el-col :span="20">
        <el-alert
          :title="alert_title"
          :type="alert_type"
          :description="alert_desc"
          @close="closeAlert"
          center
        ></el-alert>
      </el-col>
    </el-row>
    <el-row type="flex" justify="center">
      <el-col :span="24">
        <el-form
          label-position="top"
          ref="service_info"
          :model="service_info"
          class="demo-table-expand"
        >
          <el-form-item
            label="选择组件"
          >
            <el-select v-model="service_info.component_id" placeholder="请选择">
              <el-option
                v-for="item in component_list"
                :key="item.id"
                :label="item.name+':v'+item.version"
                :value="item.id"
              ></el-option>
            </el-select>
          </el-form-item>

          <el-form-item
            label="服务名"
            prop="name"
            :rule="[
            { required: true, message: '服务名不能为空'},
            ]"
          >
            <el-input type="name" v-model="service_info.name"></el-input>
          </el-form-item>
          <el-form-item
            label="服务版本"
            prop="version"
            :rule="[
            { required: true, message: '服务版本不能为空'},
            ]"
          >
            <el-input type="version" v-model="service_info.version"></el-input>
          </el-form-item>
          <el-form-item label="简介">
            <el-input v-model="service_info.desc"></el-input>
          </el-form-item>
          <el-form-item
            label="服务配置"
            prop="config"
            :rule="[
            { required: true, message: '服务配置不能为空'},
            ]"
          >
            <div>
              <p>选择配置文件上传</p>
              <input
                type="file"
                class="upload"
                @change="loadFile"
                ref="inputer"
                accept="application/json"
              />
              <el-input type="textarea" :autosize="true" v-model="service_info.config"></el-input>
              <!-- <highlight-code lang="json" :code="configString"></highlight-code> -->
              <!-- <input
              type="file"
              class="upload"
              @change="loadFile"
              ref="inputer"
              accept="application/json"
            />
              <el-button v-if="checked" type="primary" icon="el-icon-folder-add" @click="cacheConfig"></el-button>-->
            </div>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="onSubmit">立即创建</el-button>
          </el-form-item>
        </el-form>
      </el-col>
    </el-row>
  </div>
</template>

<script>
import Ajv from "ajv";
const ajv = new Ajv();
export default {
  name: "servicenew",
  data() {
    return {
      alertVisible: false,
      alert_title: "",
      alert_type: "success",
      alert_desc: "",
      component_list: [],
      service_info: {
        name: "",
        version: "",
        config: "",
        desc: "",
        component_id: ''
      }
    };
  },
  methods: {
    configObjcet: function() {
      return this.service_info.config
        ? JSON.parse(this.service_info.config)
        : {};
    },
    loadFile() {
      let inputDOM = this.$refs.inputer;
      if (inputDOM.files && inputDOM.files.length !== 0) {
        let file = inputDOM.files[0];
        let reader = new FileReader();
        reader.onload = e => {
          this.service_info.config = reader.result;
          this.alertVisible = true;
          this.alert_title = "文件读取成功";
          this.alert_type = "success";
          this.alert_desc = "";
        };
        reader.readAsText(file);
      } else {
        this.alertVisible = true;
        this.alert_title = "文件读取失败";
        this.alert_type = "error";
        this.alert_desc = "";
      }
    },
    closeAlert() {
      this.alertVisible = false;
      this.alert_title = "";
      this.alert_type = "success";
      this.alert_desc = "";
    },
    async onSubmit() {
      console.log(this.service_info)
      if (
        !(
          this.service_info.name &&
          this.service_info.version &&
          this.service_info.config &&
          this.service_info.component_id
        )
      ) {
        this.alertVisible = true;
        this.alert_title = "上传服务信息失败";
        this.alert_type = "error";
        this.alert_desc = "服务信息不满足要求";
      } else {
        let query = {
          name: this.service_info.name,
          version: this.service_info.version,
          component_id: this.service_info.component_id,
          config: this.configObjcet()
        };
        if (this.service_info.desc) {
          Object.assign(query, { desc: this.service_info.desc });
        }
        console.log(query)
        let res = await this.$axios.post(`/service/`, JSON.stringify(query), {
          headers: {
            "Content-Type": "application/json"
          }
        });
        if (res.status !== 200) {
          this.alertVisible = true;
          this.alert_title = "上传服务信息失败";
          this.alert_type = "error";
          this.alert_desc = res.data.error;
        } else {
          this.alertVisible = true;
          this.alert_title = "上传服务信息成功";
          this.alert_type = "success";
          this.alert_desc = "";
        }
      }
    }
  },
  created: function() {
    this.createdPromise = this.$axios.get(`/component/`);
  },
  mounted: async function() {
    try {
      let response = await this.createdPromise;
      this.component_list = response.data;
    } catch (error) {
      alert("数据未找到");
    }
  }
};
</script>