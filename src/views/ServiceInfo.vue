<template>
  <div class="service-info">
    <el-row type="flex" justify="center">
      <h2>服务详情</h2>

      <el-switch
        v-if="!service_online"
        v-model="service_online"
        active-color="#13ce66"
        inactive-color="#ff4949"
        active-text="线上"
        inactive-text="线下"
      ></el-switch>
      <el-switch
        v-else
        v-model="service_online"
        active-color="#13ce66"
        inactive-color="#ff4949"
        active-text="线上"
        inactive-text="线下"
        :disabled="true"
      ></el-switch>
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
          <el-form-item label="服务名">
            <el-input v-model="service_info.name" :disabled="true"></el-input>
          </el-form-item>
          <el-form-item label="服务版本">
            <el-input v-model="service_info.version" :disabled="true"></el-input>
          </el-form-item>
          <template v-if="service_info.component">
            <el-form-item label="服务所属组件ID">
              <el-input v-model="service_info.component.id" :disabled="true"></el-input>
            </el-form-item>
            <el-form-item label="服务所属组件名">
              <el-input v-model="service_info.component.name" :disabled="true"></el-input>
            </el-form-item>
            <el-form-item label="服务所属组件版本">
              <el-input v-model="service_info.component.version" :disabled="true"></el-input>
            </el-form-item>
          </template>
          <template v-if="!service_online">
            <el-form-item label="简介">
              <el-input v-model="service_info.desc"></el-input>
            </el-form-item>
            <el-form-item label="服务配置">
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
              </div>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="onSubmit">修改</el-button>
            </el-form-item>
          </template>
          <template v-else>
            <el-form-item label="简介">
              <el-input v-model="service_info.desc" :disabled="true"></el-input>
            </el-form-item>
            <el-form-item label="服务配置">
              <div>
                <el-input
                  type="textarea"
                  :autosize="true"
                  v-model="service_info.config"
                  :disabled="true"
                ></el-input>
              </div>
            </el-form-item>
          </template>
        </el-form>
      </el-col>
    </el-row>
  </div>
</template>

<script>
import Ajv from "ajv";
const ajv = new Ajv();
export default {
  name: "serviceinfo",
  data() {
    return {
      lock: true,
      service_online: true,
      service_info: {},
      alertVisible: false,
      alert_title: "",
      alert_type: "success",
      alert_desc: "",
      checked: false
    };
  },
  props: {
    id: {
      type: String,
      default: null
    }
  },
  computed: {
    // 计算属性的 getter
    configString: function() {
      return this.service_info.config
        ? JSON.stringify(this.service_info.config, null, 2)
        : "{}";
    }
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
      let Id = Number(this.id);
      if (!this.service_info.desc && !this.service_info.config) {
        this.alertVisible = true;
        this.alert_title = "更新服务信息失败";
        this.alert_type = "error";
        this.alert_desc = "服务信息不满足要求";
      } else {
        let query = {};
        if (this.service_info.desc) {
          Object.assign(query, { desc: this.service_info.desc });
        }
        if (this.service_info.config) {
          Object.assign(query, { config: this.configObjcet() });
        }
        let res = await this.$axios.put(
          `/service/${Id}`,
          JSON.stringify(query),
          {
            headers: {
              "Content-Type": "application/json"
            }
          }
        );
        if (res.status !== 200) {
          this.alertVisible = true;
          this.alert_title = "更新服务信息失败";
          this.alert_type = "error";
          this.alert_desc = res.data.error;
        } else {
          this.alertVisible = true;
          this.alert_title = "更新服务信息成功";
          this.alert_type = "success";
          this.alert_desc = "";
        }
      }
    }
  },
  watch: {
    service_online: async function(newVal, oldVal) {
      let Id = Number(this.id);
      console.log("newVal");
      console.log(newVal);
      console.log("oldVal");
      console.log(oldVal);
      if (newVal === true && oldVal == false) {
        // if (this.lock) {
        //   this.lock = false;
        //   console.log("unlocked");
        // } else {
        console.log("ok");
        let res = await this.$axios.post(`/service/${Id}/config/release`);
        if (res.status !== 200) {
          this.alertVisible = true;
          this.alert_title = "服务配置上线失败";
          this.alert_type = "error";
          this.alert_desc = res.data.error;
        } else {
          this.alertVisible = true;
          this.alert_title = "服务配置上线成功";
          this.alert_type = "success";
          this.alert_desc = "";
        }
        //}
      }
    }
  },
  created: async function() {
    let Id = Number(this.id);
    let resp = await this.$axios.get(`/service/${Id}`);
    let service_info = resp.data;
    let service_config_string = JSON.stringify(service_info.config, null, 2);
    this.service_online = service_info.online;
    this.service_info = {
      desc: service_info.desc,
      id: service_info.id,
      name: service_info.name,
      version: service_info.version,
      config: service_config_string,
      component: service_info.component
    };
  }
};
</script>