<template>
  <SecondaryLayout>
    <form>
      <div data-test="sign-in-section" class="form-section">
        <div data-test="sign-in-header-section" class="section-header">회원가입</div>
        <div data-test="sign-in-body-section" class="section-body">
          <div class="form-field">
            <div>
              <input placeholder="아이디" v-model="username" autofocus="autofocus" data-test="form-input"
                class="form-input">
            </div>
          </div>
          <div class="form-field">
            <input type="password" v-model="password" placeholder="비밀번호" data-test="form-input" class="form-input">
          </div>
          <div class="form-field">
            <input type="password" v-model="passwordVerify" placeholder="비밀번호 확인" data-test="form-input"
              class="form-input">
          </div>
          <div class="form-field">
            <input name="email" v-model="email" placeholder="이메일" data-test="form-input" class="form-input">
          </div>
          <div class="form-field">
            <input name="last_name" v-model="last_name" placeholder="성" data-test="form-input" class="form-input">
          </div>
          <div class="form-field">
            <input name="first_name" v-model="first_name" placeholder="이름" data-test="form-input" class="form-input">
          </div>
          <div class="form-field">
            <input name="age" type="number" v-model="age" placeholder="나이" data-test="form-input" class="form-input">
          </div>
          <div class="form-field">
            <!-- <label for="lang">성별</label> -->
            <select name="languages" placeholder="성별" v-model="gender" data-test="form-input" class="form-input"
              id="gender">
              <option value="">성별</option>
              <option value="F">여성</option>
              <option value="M">남성</option>
            </select>
            <!-- <input name="gender" v-model="gender" placeholder="성별" data-test="form-input" class="form-input"> -->
          </div>
          <div class="form-field">
            <input name="address" v-model="address" placeholder="주소" data-test="form-input" class="form-input">
          </div>
        </div>

        <div class="btn_login_wrap">
          <button type="submit" @click="submit" class="btn_login" id="log.login">
            <span class="btn_text">회원가입</span>
          </button>
        </div>

        <ul class="find_wrap" id="find_wrap">
          <li><a target="" href="#" class="find_text">비밀번호 찾기</a></li>
          <li><a target="" href="#" class="find_text">회원가입</a>
          </li>
        </ul>
      </div>
    </form>

    <!-- <AmplifyAuthenticator :authConfig="authConfig" ref="authenticator" /> -->
  </SecondaryLayout>
</template>

<script>
import { mapActions } from 'vuex';
// import { components } from 'aws-amplify-vue';
import { EventBus } from '@/event-bus';
import SecondaryLayout from '@/components/SecondaryLayout/SecondaryLayout';
import swal from 'sweetalert';

import { RepositoryFactory } from '@/repositories/RepositoryFactory';

const UsersRepository = RepositoryFactory.get('users');

export default {
  name: 'Auth',
  components: {
    // AmplifyAuthenticator: components.Authenticator,
    SecondaryLayout,
  },
  data() {
    return {
      username: "",
      password: "",
      passwordVerify: "",
      email: "",
      first_name: "",
      last_name: "",
      age: 19,
      gender: "",
      address: "",
      showingSignUp: undefined,
      authConfig: {
        signInConfig: {
          header: '로그인'
        },
        signUpConfig: {
          hideAllDefaults: true,
          header: 'Create account',
          signUpFields: [
            {
              label: 'Email',
              key: 'email',
              type: 'email',
              required: true
            },
            {
              label: 'Password',
              key: 'password',
              type: 'password',
              required: true
            },
            {
              label: 'Username',
              key: 'username',
              type: 'string',
              required: true
            },
          ]
        }
      }
    }
  },
  mounted() {
    if (this.showingSignUp === undefined && this.$route.query.signup) {
      this.showingSignUp = true
      EventBus.$emit('authState', 'signUp')
    } else {
      this.showingSignUp = false
    }

    this.username = "kildong.hong"
    this.password = ""
    this.passwordVerify = ""
    this.email = "kildong.hong@example.com"
    this.first_name = "길동"
    this.last_name = "홍"
    this.age = 31
    this.gender = "M"
    this.address = "테헤란로 225"
  },
  methods: {
    ...mapActions(['setVolatileUser']),
    uuidv4() {
      return ([1e7] + -1e3 + -4e3 + -8e3 + -1e11).replace(/[018]/g, c =>
        (c ^ crypto.getRandomValues(new Uint8Array(1))[0] & 15 >> c / 4).toString(16)
      );
    },
    async submit() {
      // Verify that the passwords match
      if (this.password !== this.passwordVerify) {
        swal("두 비밀번호가 일치하지 않습니다."); return;
      }

      if (this.username === "") { swal("비어있는 항목을 채워주세요."); return; }
      else if (this.password === "") { swal("비어있는 항목을 채워주세요."); return; }
      else if (this.passwordVerify === "") { swal("비어있는 항목을 채워주세요."); return; }
      else if (this.email === "") { swal("비어있는 항목을 채워주세요."); return; }
      else if (this.first_name === "") { swal("비어있는 항목을 채워주세요."); return; }
      else if (this.last_name === "") { swal("비어있는 항목을 채워주세요."); return; }
      else if (this.age < 0) { swal("비어있는 항목을 채워주세요."); return; }
      else if (this.gender === "") { swal("비어있는 항목을 채워주세요."); return; }
      else if (this.address === "") { swal("비어있는 항목을 채워주세요."); return; }

      let provisionalUserId = this.uuidv4();

      try {
        let user = {
          id: provisionalUserId,
          username: this.username,
          password: this.password,
          gender: this.gender,
          age: this.age,
          first_name: this.first_name,
          last_name: this.last_name,
          email: this.email,
          identity_id: provisionalUserId,
        }

        // console.log(user)

        const { data } = await UsersRepository.createUser(user);
        if (data.segment === "") {
          data.segment = data.gender + parseInt(data.age / 10) * 10
        }
        data['attributes'] = {}
        data['attributes']['custom:profile_user_id'] = data.username
        data['attributes']['custom:profile_email'] = data.email
        data['attributes']['custom:profile_first_name'] = data.first_name
        data['attributes']['custom:profile_last_name'] = data.last_name
        data['attributes']['custom:profile_gender'] = data.gender
        data['attributes']['custom:profile_age'] = data.age
        data['attributes']['custom:profile_segment'] = data.segment

        // console.log(data)
        this.setVolatileUser(data).then(() => {
          EventBus.$emit('authState', 'signedIn')
        })
      } catch (err) {
        // console.error(err.response.data);
        if ('response' in err) {
          swal('회원가입에 실패하였습니다.\n사유: ' + err.response.data.message);
        } else {
          swal('회원가입에 실패하였습니다.\n사유: ' + err);

        }
      }
    },
  },
  watch: {
    $route: {
      immediate: true,
      handler() {
        if (this.$route.query.signup) {
          EventBus.$emit('authState', 'signUp')
        }
      }
    }
  }
}
</script>

<style>
.form-section {
  position: relative;
  margin-bottom: 20px;
  background-color: var(--form-color);
  padding: 35px 40px;
  text-align: left;
  display: inline-block;
  min-width: 460px;
  border-radius: 6px;
  -webkit-box-shadow: 1px 1px 4px 0 rgb(0 0 0 / 15%);
  box-shadow: 1px 1px 4px 0 rgb(0 0 0 / 15%);
  -webkit-box-sizing: border-box;
  box-sizing: border-box;
}

.section-header[data-test="sign-in-header-section"] {
  color: var(--deepSquidInk);
  margin-bottom: 24px;
  font-size: 18px;
  font-weight: 500;
}

.section-body[data-test="sign-in-body-section"] {
  margin-bottom: 30px;
}

.input-label {
  color: var(--input-color);
  font-size: 14px;
  margin-bottom: 8px;
}

.form-section input[data-test="form-input"] {
  display: block;
  width: 100%;
  font-size: 14px;
  color: var(--input-color);
  background-color: var(--input-background-color);
  background-image: none;
  border: 1px solid var(--lightGrey);
  border-radius: 3px;
  -webkit-box-sizing: border-box;
  box-sizing: border-box;
  margin-bottom: 10px;
}

#gender>option {
  padding: 0px;
}

.form-section select[data-test="form-input"] {
  display: block;
  width: 100%;
  font-size: 18px;
  padding: 9px;
  color: var(--input-color);
  background-color: var(--input-background-color);
  background-image: none;
  border: 1px solid var(--lightGrey);
  border-radius: 3px;
  -webkit-box-sizing: border-box;
  box-sizing: border-box;
  margin-bottom: 10px;
}

.btn_login {
  display: block;
  width: 100%;
  padding: 13px 0 13px;
  border-radius: 6px;
  border: solid 1px rgba(0, 0, 0, .15);
  color: var(--button-color);
  background-color: var(--button-background-color);
  box-sizing: border-box;
}

button[data-test="sign-in-sign-in-button"] {
  min-width: 153px;
  display: inline-block;
  margin-bottom: 0;
  font-size: 12px;
  font-weight: 400;
  line-height: 1.42857143;
  text-align: center;
  white-space: nowrap;
  vertical-align: middle;
  -ms-touch-action: manipulation;
  touch-action: manipulation;
  cursor: pointer;
  -webkit-user-select: none;
  -moz-user-select: none;
  -ms-user-select: none;
  user-select: none;
  background-image: none;
  color: var(--button-color);
  background-color: var(--button-background-color);
  border-color: #ccc;
  text-transform: uppercase;
  padding: 14px 0;
  letter-spacing: 1.1px;
  border: none;
  --webkit-box-direction: reverse;
}

.find_wrap,
.relogin_find_wrap {
  padding: 1em 0em;
  text-align: center;

  margin-inline-start: 0px;
  margin-inline-end: 0px;
}

.find_wrap li,
.relogin_find_wrap li {
  position: relative;
  display: inline-block;
  margin: 0;
  padding: 0em 0.5em;
  -webkit-text-size-adjust: none;
  font-family: -apple-system, BlinkMacSystemFont, helvetica, "Apple SD Gothic Neo", sans-serif;
}

.find_wrap .en_find_text,
.find_wrap .find_text,
.relogin_find_wrap .en_find_text,
.relogin_find_wrap .find_text {
  display: inline-block;
  font-size: 14px;
  line-height: 17px;
  text-decoration: none;
  color: #888;
}

/* Amplify Auth Form Styling */

div[data-test="sign-in-section"],
div[data-test="sign-up-section"],
div[data-test="verify-contact-section"],
div[data-test="require-new-password-section"],
div[data-test="federated-sign-in-section"],
div[data-test="confirm-sign-up-section"],
div[data-test="confirm-sign-in-section"],
div[data-test="set-mfa-section"],
div[data-test="forgot-password-section"] {
  box-shadow: none;
}

/* On xs screens, override default min-width on forms to prevent overflow */
@media (max-width: 576px) {

  div[data-test="sign-in-section"],
  div[data-test="sign-up-section"],
  div[data-test="verify-contact-section"],
  div[data-test="require-new-password-section"],
  div[data-test="federated-sign-in-section"],
  div[data-test="confirm-sign-up-section"],
  div[data-test="confirm-sign-in-section"],
  div[data-test="set-mfa-section"],
  div[data-test="forgot-password-section"] {
    min-width: initial !important;
  }
}

/* Set font-size to 18px to disable auto-zoom on mobile Safari */
div[data-test="sign-in-section"] input,
div[data-test="sign-up-section"] input,
div[data-test="verify-contact-section"] input,
div[data-test="require-new-password-section"] input,
div[data-test="federated-sign-in-section"] input,
div[data-test="confirm-sign-up-section"] input,
div[data-test="confirm-sign-in-section"] input,
div[data-test="set-mfa-section"] input,
div[data-test="forgot-password-section"] input {
  font-size: 18px !important;
  padding: .5em;
}

/* Make links in form text/labels more noticeable */
div[data-test="sign-in-section"] a,
div[data-test="sign-up-section"] a,
div[data-test="verify-contact-section"] a,
div[data-test="require-new-password-section"] a,
div[data-test="federated-sign-in-section"] a,
div[data-test="confirm-sign-up-section"] a,
div[data-test="confirm-sign-in-section"] a,
div[data-test="set-mfa-section"] a,
div[data-test="forgot-password-section"] a {
  color: var(--blue-500) !important;
}

/* Make error messages stand out and match bootstrap alert-error */
div[data-test="sign-in-section"] div.error,
div[data-test="sign-up-section"] div.error,
div[data-test="verify-contact-section"] div.error,
div[data-test="require-new-password-section"] div.error,
div[data-test="federated-sign-in-section"] div.error,
div[data-test="confirm-sign-up-section"] div.error,
div[data-test="confirm-sign-in-section"] div.error,
div[data-test="set-mfa-section"] div.error,
div[data-test="forgot-password-section"] div.error {
  color: #721c24;
  background-color: #f8d7da;
  position: relative;
  padding: .75rem 1.25rem;
  margin-top: 1rem;
  border: 1px solid #f5c6cb;
  border-radius: .25rem;
}
</style>
