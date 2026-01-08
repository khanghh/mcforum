export function userHasRole(user, role) {
  return user && user.role == role
}

export function userHasAnyRole(user, ...roles) {
  if (!roles || !roles.length) {
    return false
  }
  for (let i = 0; i < roles.length; i++) {
    if (userHasRole(user, roles[i])) {
      return true
    }
  }
  return false
}

export function userIsOwner(user) {
  return userHasRole(user, 'owner')
}

export function userIsAdmin(user) {
  return userHasRole(user, 'admin')
}

export function userIsManager(user) {
  return userHasAnyRole(user, 'owner', 'admin', 'moderator')
}