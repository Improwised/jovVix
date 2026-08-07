local claims = {
  email_verified: false,
} + std.extVar('claims');

local cap(value, maxLength) =
  if std.length(value) > maxLength then std.substr(value, 0, maxLength) else value;

{
    identity: {
        traits: {
            [if 'email' in claims && claims.email_verified then 'email' else null]: claims.email,
            name: {
                first: if 'given_name' in claims then cap(claims.given_name, 12) else null,
                last: if 'family_name' in claims then cap(claims.family_name, 50) else null,
            },
        },
    },
}
