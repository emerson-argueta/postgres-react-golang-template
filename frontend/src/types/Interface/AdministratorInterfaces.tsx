import { IChurch } from "./ChurchInterfaces";

enum SubscriptionType {
  FreePlan = "free",
  BasicPlan = "basic",
  StandardPlan = "standard",
  PremiumPlan = "premium",
}

export enum Access {
  NonRestricted = 'non-restricted',
  Restricted = 'restricted'
}
export enum Role {
  Creator = 'creator',
  Support = 'support'
}
export const AccessTypes = {
  [Access.Restricted.toString()]: Access.Restricted.toString(),
  [Access.NonRestricted.toString()]: Access.NonRestricted.toString()
}
export const RoleTypes = {
  [Role.Creator.toString()]: Role.Creator.toString(),
  [Role.Support.toString()]: Role.Support.toString(),
}
export interface IChurches {
  [id: number]: IChurch
}
export interface IPaymentGateway {
  [index: string]: any
}
export interface ISubscription {
  freeusagelimitcount: number,
  customeremail: string
  type: SubscriptionType,
  paymentgateway: IPaymentGateway
}
export interface IToken {
  accesstoken?: string,
  refreshtoken?: string
}
export interface IAdministrator {
  uuid?: string,
  firstname?: string,
  lastname?: string,
  address?: string,
  phone?: string,
  churches?: IChurches,
  subscription?: ISubscription,
  email?: string,
  password?: string,
  token?: IToken,
  access?: Access,
  role?: Role

}