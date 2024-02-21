<template>
  <Body>
    <header class="container-fluid ps-0 pe-0">
      <nav class="navbar navbar-expand-lg navbar-dark bg-dark">
        <div class="container-fluid p-2">
          <!-- Title on the left -->
          <a class="navbar-brand" href="#">Quiz App</a>

          <!-- Profile/Login on the right -->
          <div class="navbar-collapse justify-content-end">
            <ul class="navbar-nav">
              <li class="nav-item">
                <a class="nav-link" href="#">Profile/Login</a>
              </li>
            </ul>
          </div>
        </div>
      </nav>
    </header>
    <main class="container">
      <slot />
    </main>
    <footer
      class="container-fluid position-fixed bottom-0 bg-dark text-capitalize fs-5 p-2"
    >
      <div class="text-center p-3 text-light">
        © 2024 Copyright:
        <a class="text-light" href="/">Quizz_app</a>
      </div>
    </footer>
    <footer class="bg-body-tertiary text-center text-lg-start"></footer>
  </Body>
</template>
