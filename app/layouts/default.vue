<template>
  <div
    class="container-fluid d-flex align-items-center min-vh-100 position-relative"
  >
    <header class="d-flex bg-dark header position-fixed start-0 end-0 top-0">
      <nav class="navbar navbar-expand-md navbar-dark w-100">
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
      class="container-md d-flex justify-content-center align-items-center h-100 pageContainer px-0"
    >
      <slot />
    </main>
    <footer
      class="d-flex justify-content-center bg-dark position-absolute start-0 end-0 bottom-0"
    >
      <div class="text-center text-capitalize p-3 text-light w-100">
        © {{ new Date().getFullYear() }} Copyright:
        <NuxtLink class="text-light" to="/">Quizz_app</NuxtLink>
      </div>
    </footer>
  </div>
</template>
<style scoped>
header,
footer {
  z-index: 10;
}
.pageContainer {
  padding-top: 100px;
  padding-bottom: 100px;
}
</style>
