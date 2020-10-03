import { IAccountStatement } from "./Interfaces";
import { IDonation } from "./DonationInterfaces";

export interface IChurches {
  donationcount?: number,
  firstdonation?: Date,
}
export type IDonations = Array<IDonation>

export interface IDonator {
  id?: number,
  uuid?: string,
  firstname?: string,
  lastname?: string,
  email?: string,
  address?: string,
  phone?: string,
  churches?: IChurches,
  accountstatement?: IAccountStatement,
  donations?: IDonations
}