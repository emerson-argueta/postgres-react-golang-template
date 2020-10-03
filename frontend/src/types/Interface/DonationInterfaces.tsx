export enum DonationType {
    Online = 'online',
    Cash = 'cash',
    Check = 'check'
}
export const DonationTypes = {
    [DonationType.Online]:DonationType.Online,
    [DonationType.Cash]:DonationType.Cash,
    [DonationType.Check]:DonationType.Check,
    
}

export interface IDonation {
    donatorid?:number,
    churchid?:number,
    date?: string,
    amount?: number,
    type?: DonationType,
    currency?: string,
    category?: string,
    account?:string,
    details?: string
}