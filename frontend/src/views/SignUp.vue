<template>
  <div class="main">
    <div class="container">
      <h2 class="form-title">Sign Up</h2>
      <div class="form-group">
        <label for="name"><span style="color:red;">* </span>Username</label>
        <el-input type="text" required name="name" id="name" placeholder="Username" v-model="username" />
      </div>
      <div class="form-group">
        <label for="email"><span style="color:red;">* </span>Email</label>
        <el-input type="email" required name="email" id="email" placeholder="Email" v-model="email" />
      </div>
      <div class="form-group">
        <label for="pass"><span style="color:red;">* </span>Password</label>
        <el-input type="password" required name="pass" id="pass" placeholder="Password" v-model="password" />
      </div>
      <div class="form-group">
        <label for="re_pass"><span style="color:red;">* </span>Comfirm Password</label>
        <el-input type="password" required name="re_pass" id="re_pass" placeholder="Re-enter Password" v-model="confirm_password" />
      </div>
      <div class="form-group">
        <label for="gender"><span style="color:red;">* </span>Gender</label>
        <div id="gender">
          <el-radio v-model="gender" :label="1">Male</el-radio>
          <el-radio v-model="gender" :label="2">Female</el-radio>
        </div>
      </div>
      <div class="form-btn">
        <button type="button" class="btn btn-info" @click="submit">Submit</button>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "SignUp",
  data() {
    return {
      username: "",
      password: "",
      email: '',
      gender: 1,
      confirm_password: "",
      submitted: false
    };
  },
  computed: {
  },
  created() {

  },
  methods: {
    submit() {
      this.$axios({
        method: 'post',
        url: '/signup',
        data: {
          username: this.username,
          email: this.email,
          gender: this.gender,
          password: this.password,
          confirm_password: this.confirm_password
        }
      }).then((res) => {
        console.log(res.data);
        if (res.code == 1000) {
          console.log('signup success');
          this.$router.push({ name: "Login" });
        } else {
          console.log(res.msg);
        }
      }).catch((error) => {
        console.log(error)
      })
    }
  }
};
</script>
<style lang="less" scoped>
.main {
  background: #6190E8;
  /* fallback for old browsers */
  background: -webkit-linear-gradient(to right, #A7BFE8, #6190E8);
  /* Chrome 10-25, Safari 5.1-6 */
  background: linear-gradient(to right, #A7BFE8, #6190E8);
  /* W3C, IE 10+/ Edge, Firefox 16+, Chrome 26+, Opera 12+, Safari 7+ */

  padding: 150px 0;
  min-height: 60vh;

  .container {
    width: 600px;
    background: #fff;
    margin: 0 auto;
    max-width: 1200px;
    padding: 20px;

    .form-title {
      margin-bottom: 33px;
      text-align: center;
    }

    .form-group {
      margin: 15px;

      label {
        display: inline-block;
        max-width: 100%;
        margin-bottom: 5px;
        font-weight: 700;
      }

      // .form-control {
      //   display: block;
      //   width: 100%;
      //   height: 34px;
      //   padding: 6px 12px;
      //   font-size: 14px;
      //   line-height: 1.42857143;
      //   color: #555;
      //   background-color: #fff;
      //   background-image: none;
      //   border: 1px solid #ccc;
      //   border-radius: 4px;
      // }
    }

    .form-btn {
      display: flex;
      justify-content: center;

      .btn {
        padding: 6px 20px;
        font-size: 18px;
        line-height: 1.3333333;
        border-radius: 6px;
        display: inline-block;
        margin-bottom: 0;
        font-weight: 400;
        text-align: center;
        white-space: nowrap;
        vertical-align: middle;
        -ms-touch-action: manipulation;
        touch-action: manipulation;
        cursor: pointer;
        border: 1px solid transparent;
      }

      .btn-info {
        color: #fff;
        background-color: #5bc0de;
        border-color: #46b8da;
      }
    }
  }
}
</style>