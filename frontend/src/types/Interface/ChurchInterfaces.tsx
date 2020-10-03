import { IDonator } from "./DonatorInterfaces";
import { IAccountStatement } from "./Interfaces";
import { IAdministrator, Access, Role } from "./AdministratorInterfaces";

export interface IAdministrators {
    [uuid: string]: IAdministrator
}
export interface IDonators {
    [id: number]: IDonator
}

export interface IChurch {
    id?: number,
    type?: string,
    name?: string,
    address?: string,
    phone?: string,
    email?: string,
    password?: string,
    administrators: IAdministrators,
    donators?: IDonators,
    accountstatement?: IAccountStatement,
    donationcategories?: { [index: string]: string },
    churchdonationreport: IChurchDonationReport,
    donationstatementreport: IDonationStatementReport,
    access: Access,
    role: Role

}
export interface IChurchDonationReport {
    churchid: number,
    reporttype: string,
    donationcategories: Array<string>
    donations?: { [donatorid: number]: any }
    donationssum?: { [category: string]: any }
    sumfilter: { timeperiod: string, multiplier: number }
    timerange: { upper: string, lower: string }
}
export interface IDonationStatementReport {
    churchid: number,
    donatorids: Array<number>,
    donations?: { [donatorid: number]: any }
    donationssum?: { [category: string]: any }
    timerange: { upper: string, lower: string }
}
export const SumFilter = {
    day: "day",
    week: "week",
    month: "month",
    year: "year"
}

export const ReportType = {
    donations: "donations"
    // values can be added in the future
}