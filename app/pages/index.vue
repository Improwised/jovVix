<template>
  <div
    class="full-height-layout d-flex align-items-center justify-content-center px-15px py-10"
  >
    <div class="row align-items-center">
      <div class="col-md-7">
        <div class="quiz-content text-md-start text-center">
          <h1 class="text-72px font-weight-bold text-black mb-5">
            Welcome to the Quiz App
          </h1>
          <p class="text-24px text-black mb-5">Explore, Learn, Quiz!</p>
        </div>
      </div>
      <div class="col-md-5">
        <lottie-player
          autoplay
          loop
          mode="normal"
          src="testimg.json"
          style="width: 100%"
        ></lottie-player>
        <!-- vector image -->
        <img
          src="../assets/images/shapes/1.svg"
          alt=""
          class="shapes shape-one absolute z--1 left-20% top-9%"
        />
        <!-- Purple Half Circle -->
        <img
          src="../assets/images/shapes/2.svg"
          alt=""
          class="shapes shape-three absolute z--1 left-11% top-35%"
        />
        <!-- Birds -->
        <img
          src="../assets/images/shapes/3.svg"
          alt=""
          class="shapes shape-four absolute z--1 left-10% top-70%"
        />
        <!-- cyan circle -->
        <img
          src="../assets/images/shapes/4.svg"
          alt=""
          class="shapes shape-five absolute z--1 left-22% top-90%"
        />
        <!-- > sign -->
        <img
          src="../assets/images/shapes/5.svg"
          alt=""
          class="shapes shape-six absolute z--1 right-22% top-9%"
        />
        <!-- < sign -->
        <img
          src="../assets/images/shapes/6.svg"
          alt=""
          class="shapes shape-seven absolute z--1 right-12% top-28%"
        />
        <!-- pink circle -->
        <img
          src="../assets/images/shapes/7.svg"
          alt=""
          class="shapes shape-eight absolute z--1 right-12% top-60%"
        />
        <!-- line -->
        <img
          src="../assets/images/shapes/8.svg"
          alt=""
          class="shapes shape-nine absolute z--1 right-1/4 top-85%"
        />
        <!-- small circle -->
        <img
          src="../assets/images/shapes/9.svg"
          alt=""
          class="shapes shape-ten absolute z--1 right-5% bottom-0 md:hidden sm:hidden xsm:hidden"
        />
        <!-- half yellow circle -->
        <img
          src="../assets/images/shapes/8.svg"
          alt=""
          class="shapes shape-eleven absolute z--1 right-1/4 top-85%"
        />
        <img
          src="../assets/images/shapes/8.svg"
          alt=""
          class="shapes shape-twelve absolute z--1 right-1/4 top-85%"
        />
        <img
          src="../assets/images/shapes/8.svg"
          alt=""
          class="shapes shape-thirteen absolute z--1 right-1/4 top-85%"
        />
        <img
          src="../assets/images/shapes/2.svg"
          alt=""
          class="shapes shape-fourteen absolute z--1 left-11% top-35%"
        />
        <img
          src="../assets/images/boxes.png"
          alt=""
          class="shapes img-boxes absolute z--1 left-11% top-35%"
        />
        <!-- 4 box image -->
        <img
          src="../assets/images/shapes/3.svg"
          alt=""
          class="shapes shape-fifteen absolute z--1 left-10% top-70%"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { useUsersStore } from "~~/store/users";
const userData = useUsersStore();
const { setUserData } = userData;
const { api_url } = useRuntimeConfig().public;
const headers = useRequestHeaders(["cookie"]);

const setUserDataStore = async () => {
  try {
    const response = await fetch(api_url + "/user/who", {
      method: "GET",
      credentials: "include",
      headers: headers,
      mode: "cors",
    });
    if (response.status != 200) {
      throw new Error(response.status);
    } else if (response.status == 200) {
      const data = await response.json();
      setUserData(data?.data?.role);
    }
  } catch (error) {
    if (error.message == 401) {
      console.log(error.message);
      setUserData(null);
      return;
    }
    console.log(error.message);
  }
};

onMounted(async () => {
  import("@lottiefiles/lottie-player");
  await setUserDataStore();
});
</script>

<style scoped>
.homepage {
  font-size: 30px;
}

@keyframes moveUpDown {
  0%,
  100% {
    transform: translateY(0);
  }

  50% {
    transform: translateY(-20px);
  }
}

@keyframes opacitySwitch {
  0%,
  100% {
    opacity: 1;
  }

  50% {
    opacity: 0.5;
  }
}

.shape-one {
  position: absolute;
  left: 0;
  top: 5%;
  animation: opacitySwitch 1.5s infinite;
}

.shape-two {
  position: absolute;
  left: 20%;
  top: 9%;
  animation: opacitySwitch 1.5s infinite;
}

.shape-three {
  position: absolute;
  left: 6%;
  top: 55%;
  animation: moveUpDown 2s infinite ease-in-out;
}

.shape-four {
  position: absolute;
  left: 10%;
  top: 70%;
  animation: opacitySwitch 1.5s infinite;
}

.shape-five {
  position: absolute;
  left: 22%;
  top: 90%;
  animation: moveUpDown 2s infinite ease-in-out;
}

.shape-six {
  position: absolute;
  right: 22%;
  top: 9%;
  animation: moveUpDown 2s infinite ease-in-out;
}

.shape-seven {
  position: absolute;
  right: 12%;
  top: 78%;
  animation: opacitySwitch 1.5s infinite;
}

.shape-eight {
  position: absolute;
  right: 12%;
  top: 60%;
  animation: moveUpDown 2s infinite ease-in-out;
}

.shape-nine {
  position: absolute;
  right: 25%;
  top: 85%;
  animation: moveUpDown 2s infinite ease-in-out;
}

.shape-ten {
  position: absolute;
  right: 5%;
  bottom: 0;
  animation: opacitySwitch 1.5s infinite;
}

.shape-eleven {
  position: absolute;
  right: 55%;
  bottom: 71%;
  animation: moveUpDown 2s infinite ease-in-out;
}

.shape-twelve {
  position: absolute;
  right: 15%;
  bottom: 50%;
  animation: moveUpDown 2s infinite ease-in-out;
}

.shape-thirteen {
  position: absolute;
  right: 33%;
  bottom: 0;
  animation: moveUpDown 2s infinite ease-in-out;
}

.shape-fourteen {
  position: absolute;
  left: 51%;
  top: 95%;
  animation: moveUpDown 2s infinite ease-in-out;
}

.shape-fifteen {
  position: absolute;
  left: 10%;
  top: 70%;
  animation: opacitySwitch 1.5s infinite;
}

.img-boxes {
  position: absolute;
  left: 21%;
  top: 12%;
  animation: moveUpDown 2s infinite ease-in-out;
}

.landing-vector-img {
  max-width: 100%;
  /* Ensure image doesn't overflow its container */
  height: auto;
  /* Maintain aspect ratio */
}

.quiz-content {
  font-family: "Gordita", sans-serif;
  max-width: 860px;
  text-align: center;
  margin: 0 auto;
  padding-top: 10%;
}

.quiz-content h1 {
  font-size: 60px;
  font-weight: bold;
  line-height: 1.18em;
  color: black;
  padding: 35px 100 60px;
  padding-top: 10%;
}

.quiz-content p {
  font-size: 24px;
  line-height: 1.79em;
  color: black;
  padding-bottom: 5px;
}

@media (max-width: 768px) {
  .shape-ten {
    display: none;
  }

  .quiz-content p {
    padding-bottom: -50px;
  }
}
</style>
