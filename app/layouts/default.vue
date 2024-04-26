<template>
  <div class="vh-100">
    <header
      class="container-fluid ps-0 pe-0 position-absolute top-0 left-0 right-0"
    >
      <nav class="navbar navbar-expand-lg navbar-dark bg-dark">
        <div class="container-fluid p-2">
          <!-- Title on the left -->
          <NuxtLink class="navbar-brand" to="/">Quiz App</NuxtLink>
          <button
            class="navbar-toggler"
            type="button"
            data-bs-toggle="collapse"
            data-bs-target="#navbarNavAltMarkup"
            aria-controls="navbarNavAltMarkup"
            aria-expanded="false"
            aria-label="Toggle navigation"
          >
            <span class="navbar-toggler-icon"></span>
          </button>

          <!-- Profile/Login on the right -->
          <div
            id="navbarNavAltMarkup"
            class="collapse navbar-collapse justify-content-end"
          >
            <ul class="navbar-nav">
              <li class="nav-item">
                <NuxtLink class="nav-link" to="/account/login"
                  >Profile/Login</NuxtLink
                >
              </li>
              <li>
                <NuxtLink class="nav-link" to="/join">Join</NuxtLink>
              </li>
              <li>
                <NuxtLink class="nav-link" to="/admin/quiz/list-quiz"
                  >Admin</NuxtLink
                >
              </li>
            </ul>
          </div>
        </div>
      </nav>
    </header>
    <main
      class="container d-flex justify-content-center align-items-center h-100"
    >
      <slot />
    </main>
    <footer
      class="container-fluid position-fixed bottom-0 bg-dark text-capitalize fs-5 p-2"
    >
      <div class="text-center p-3 text-light">
        © {{ new Date().getFullYear() }} Copyright:
        <NuxtLink class="text-light" to="/">Quizz_app</NuxtLink>
      </div>
    </footer>
    <footer class="bg-body-tertiary text-center text-lg-start"></footer>
  </div>
</template>
<style scoped>
header,
footer {
  z-index: 10;
}
</style>
