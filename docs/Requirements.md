# Project Rename Requirements: quizz-app → Jovvix

## Overview
This document outlines the requirements and changes needed to complete the rename of the project from "quizz-app" to "Jovvix". Jovvix is an online, multi-user rapid quiz web application.

## Completed Changes ✅

### Documentation Files
- [x] **README.md** - Updated project title, descriptions, and installation commands
- [x] **CONTRIBUTING.md** - Updated project references
- [x] **CODE_OF_CONDUCT.md** - Updated community references

### Frontend Application Files
- [x] **app/nuxt.config.ts** - Updated page title from "Quiz App" to "Jovvix"
- [x] **app/pages/recovery.vue** - Updated brand name in recovery page
- [x] **app/pages/index.vue** - Updated welcome message
- [x] **app/pages/account/login.vue** - Updated brand name in login page
- [x] **app/components/Header.vue** - Updated header comment

### Configuration Files (Already Updated)
- [x] **docker-compose.yaml** - Container names, environment variables, and database names
- [x] **api/docker-compose.yaml** - Container names and database configurations
- [x] **api/Dockerfile** - Working directory and binary names
- [x] **api/.air.toml** - Build commands and binary paths

## Verification Requirements

### 1. Functional Testing
- [ ] Verify application builds successfully with Docker
- [ ] Test that all services start correctly
- [ ] Confirm database connections work with new naming
- [ ] Validate frontend loads and displays "Jovvix" branding

### 2. Documentation Review
- [ ] Ensure all user-facing documentation reflects "Jovvix" branding
- [ ] Verify installation instructions work with new project name
- [ ] Check that all internal references are consistent

### 3. Configuration Validation
- [ ] Confirm Docker containers start with correct names
- [ ] Validate environment variables are properly set
- [ ] Test database migrations run successfully
- [ ] Verify API endpoints respond correctly

## Additional Considerations

### SEO and Branding
- Consider updating meta descriptions and keywords if applicable
- Update any social media or external references
- Review analytics configurations for new branding

### Repository Settings
- Update repository description on GitHub
- Consider updating repository topics/tags
- Update any CI/CD pipeline names or references

### Future Maintenance
- Monitor for any missed references in logs or error messages
- Update any external documentation or wikis
- Consider creating redirects if the old name was used in URLs

## Testing Checklist

### Local Development
- [ ] `docker-compose up --build` completes successfully
- [ ] Frontend accessible at http://127.0.0.1:5000
- [ ] API accessible at http://127.0.0.1:3000
- [ ] Database connections established
- [ ] All services communicate properly

### User Interface
- [ ] Page titles show "Jovvix"
- [ ] Navigation elements display correct branding
- [ ] Login/registration pages show updated name
- [ ] Error pages maintain consistent branding

### API and Backend
- [ ] Swagger documentation accessible
- [ ] Database migrations complete without errors
- [ ] Authentication flows work correctly
- [ ] WebSocket connections establish properly

## Rollback Plan

If issues arise during deployment:
1. Revert to previous Docker images
2. Restore database with original naming if needed
3. Update DNS/routing back to original configuration
4. Communicate changes to users if necessary

## Success Criteria

The rename is considered complete when:
- All user-facing elements display "Jovvix" consistently
- Application functions identically to pre-rename state
- No broken links or references to old name exist
- Documentation accurately reflects new branding
- All automated tests pass with new configuration

## Notes

- The rename primarily affects branding and display names
- Core functionality and architecture remain unchanged
- Database schema and API endpoints maintain compatibility
- Existing user data and configurations are preserved
