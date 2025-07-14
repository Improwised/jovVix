<template>
  <div>
    <div
      class="full-height-layout d-flex align-items-center justify-content-center px-15px py-10"
    >
      <div class="row align-items-center">
        <div class="col-md-7">
          <div class="text-md-start text-center">
            <h1 class="content-font text-72px text-black mb-5">
              Welcome to the jovVix
            </h1>
            <p class="text-24px text-black mb-5">Explore, Learn, Quiz!</p>
          </div>
        </div>
        <div class="col-md-5">
          <lottie-player
            autoplay
            loop
            speed="0.5"
            mode="normal"
            src="landing_page_player.json"
            style="width: 100%"
          ></lottie-player>
          <!-- vector image -->

          <img
            src="../assets/images/shapes/2.svg"
            alt=""
            class="shapes shape-three absolute z--1 left-11% top-35%"
          />
          <!-- cyan circle -->
          <img
            src="../assets/images/shapes/4.svg"
            alt=""
            class="shapes shape-five absolute z--1 left-22% top-90%"
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
        </div>
      </div>
    </div>

    <!-- Informational Banner -->
    <div class="container py-5">
      <div class="row justify-content-center">
        <div class="col-lg-10">
          <div class="banner-card">
            <div class="row align-items-center">
              <!-- Icon Section -->
              <div class="col-md-2 text-center mb-3 mb-md-0">
                <div class="banner-icon">
                  <font-awesome-icon :icon="['fas', 'info-circle']" class="fs-1 text-primary" />
                </div>
              </div>

              <!-- Content Section -->
              <div class="col-md-10">
                <div class="banner-content">
                  <div class="d-flex align-items-center mb-3">
                    <span class="badge bg-primary me-2">DEMO</span>
                    <h4 class="mb-0 fw-bold text-dark">Public Demo Instance</h4>
                  </div>

                  <div class="alert alert-info border-0 mb-4" role="alert">

                    <div>
                      <strong>Info:</strong> This demo supports up to <span class="fw-bold text-primary">100 concurrent
                        users</span>.
                      Performance may be affected beyond this limit.
                    </div>
                  </div>
                </div>

                <div class="row g-4">
                  <!-- Deploy Section -->
                  <div class="col-md-6">
                    <div class="feature-card h-100">
                      <div class="d-flex align-items-center mb-2">
                        <font-awesome-icon :icon="['fas', 'rocket']" class="me-2 text-success" />
                        <h6 class="mb-0 fw-semibold">Deploy Your Own</h6>
                      </div>
                      <p class="text-muted mb-3 small">
                        Get unlimited users and full control with your own instance
                      </p>
                      <!-- <button class="btn btn-outline-success btn-sm">
                            <font-awesome-icon :icon="['fas', 'book']" class="me-1" />
                            View Guide
                          </button> -->
                    </div>
                  </div>

                  <!-- Contact Section -->
                  <div class="col-md-6">
                    <div class="feature-card h-100">
                      <div class="d-flex align-items-center mb-2">
                        <font-awesome-icon :icon="['fas', 'headset']" class="me-2 text-info" />
                        <h6 class="mb-0 fw-semibold">Get in touch:</h6>
                      </div>
                      <div class="d-flex flex-wrap gap-2">
                        <a href="mailto:contact@jovvix.com"
                          class="btn btn-outline-primary btn-sm text-decoration-none contact-btn"
                        >
                          <font-awesome-icon :icon="['fas', 'envelope']" class="me-1" />
                          contact@jovvix.com
                        </a>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
onMounted(async () => {
  import("@lottiefiles/lottie-player");

  // Initialize Bootstrap tooltips
  if (import.meta.client) {
    const { Tooltip } = await import('bootstrap');
    const tooltipTriggerList = document.querySelectorAll('[data-bs-toggle="tooltip"]');
    [...tooltipTriggerList].map(tooltipTriggerEl => new Tooltip(tooltipTriggerEl));
  }
});
</script>

<style scoped>
.homepage {
  font-size: 30px;
}

.demo-info-banner {
  background: linear-gradient(135deg, #f8f9fa 0%, #e9ecef 100%);
  border-top: 1px solid #dee2e6;
  position: relative;
  overflow: hidden;
}

.demo-info-banner::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 4px;
  background: linear-gradient(90deg, #007bff 0%, #28a745 50%, #17a2b8 100%);
}

.banner-card {
  background: white;
  border-radius: 16px;
  padding: 2rem;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  backdrop-filter: blur(10px);
  transition: transform 0.3s ease, box-shadow 0.3s ease;
}

.banner-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.15);
}

.banner-icon {
  background: linear-gradient(135deg, #e3f2fd 0%, #bbdefb 100%);
  border-radius: 50%;
  width: 80px;
  height: 80px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto;
  box-shadow: 0 4px 16px rgba(0, 123, 255, 0.2);
}

.feature-card {
  background: linear-gradient(135deg, #f8f9fa 0%, #ffffff 100%);
  border: 1px solid #e9ecef;
  border-radius: 12px;
  padding: 1.5rem;
  transition: all 0.3s ease;
  position: relative;
  overflow: hidden;
}

.feature-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 3px;
  background: linear-gradient(90deg, transparent 0%, #007bff 50%, transparent 100%);
  transform: translateX(-100%);
  transition: transform 0.3s ease;
}

.feature-card:hover::before {
  transform: translateX(100%);
}

.feature-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.1);
  border-color: #007bff;
}

.btn {
  border-radius: 8px;
  font-weight: 500;
  transition: all 0.3s ease;
  text-transform: none;
}

.btn:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.contact-btn {
  position: relative;
  overflow: hidden;
}

.contact-btn::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.2), transparent);
  transition: left 0.5s;
}

.contact-btn:hover::before {
  left: 100%;
}

/* Custom tooltip styling */
.tooltip {
  font-size: 0.875rem;
}

.tooltip .tooltip-inner {
  background-color: #212529;
  color: white;
  border-radius: 6px;
  padding: 8px 12px;
  font-weight: 500;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.tooltip .tooltip-arrow::before {
  border-top-color: #212529;
}

.badge {
  font-size: 0.75rem;
  padding: 0.5rem 0.75rem;
  border-radius: 6px;
  font-weight: 600;
}

@media (max-width: 768px) {
  .banner-card {
    padding: 1.5rem;
  }

  .feature-card {
    padding: 1rem;
  }

  .banner-icon {
    width: 60px;
    height: 60px;
  }
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
