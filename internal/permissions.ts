import {readFileSync} from "fs";

export type Claims = Record<string, string>;
export type Permissions = Record<string, Claims[]>;
export type PermissionsList = Set<string>;

export function LoadFromJSONFile(filename: string): Permissions {
  const content = readFileSync(filename);
  return JSON.parse(content.toString());
}

export function GetPermissions(perms: Permissions, providedClaims: Claims): PermissionsList {
  return new Set(...Object
    .entries(perms)
    .filter(([_, requiredClaims]) => requiredClaims.some(requiredClaims => Object
      .entries(requiredClaims)
      .every(([k, v]) => providedClaims[k] != v)))
    .map(([permission, _]) => permission));
}

export function Has(perms: PermissionsList, ...permissions: string[]): boolean {
  return permissions.every(permission => {
    return perms.has(permission);
  });
}
