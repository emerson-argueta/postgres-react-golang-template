import React, { Dispatch } from 'react'
import MaterialTable, { Column, MaterialTableProps } from 'material-table'
import * as DONATOR_TYPES from '../../../types/Interface/DonatorInterfaces'
import * as DONATION_TYPES from '../../../types/Interface/DonationInterfaces'
import * as CHURCH_TYPES from '../../../types/Interface/ChurchInterfaces'
import { useDispatch } from 'react-redux'
import { updateDonationACT } from '../../../redux/actions/AppActions'
import { IconButton, Tooltip } from '@material-ui/core'
import GetAppIcon from '@material-ui/icons/GetApp'
import jsPDF from 'jspdf';
import autoTable, { applyPlugin, CellInput, ColumnInput } from 'jspdf-autotable';

applyPlugin(jsPDF)

interface StaticDynamicColumns extends Column<any> {
    static: boolean
}

export const DonationsTable = (
    { church, donators, donationReport, statementReport }:
        {
            church: CHURCH_TYPES.IChurch,
            donators: CHURCH_TYPES.IDonators | null,
            donationReport?: CHURCH_TYPES.IChurchDonationReport | null,
            statementReport?: CHURCH_TYPES.IDonationStatementReport | null

        }
) => {
    const dispatch = useDispatch()

    let materialTables: Array<JSX.Element> = []
    let statementMaterialTables: Array<JSX.Element> = []
    let donatorArray = Object.values(donators || {})
    if (donatorArray.length !== 0 && !donationReport && !statementReport) {

        materialTables = prepareDonationsTable(donatorArray, church, dispatch)

    } else if (donationReport && donationReport.donations && donationReport.donationssum && donators) {

        materialTables = prepareChurchDonationReportTables(donationReport, church, donators)

    } else if (statementReport && statementReport.donations && statementReport.donationssum && donators) {

        statementMaterialTables = prepareDonationStatementReportTables(statementReport, church, donators)

    }

    return (
        <div>
            <Tooltip title={'Download All'}>
                <IconButton
                    color="inherit"
                    onClick={(event) =>
                        exportAllTablesToPDF(church, materialTables.length !== 0 ? materialTables : statementMaterialTables)
                    }

                >
                    <GetAppIcon />
                </IconButton>
            </Tooltip>
            {materialTables}
        </div>
    )
}

const prepareDonationsTable = (
    tableDonators: Array<DONATOR_TYPES.IDonator>,
    church: CHURCH_TYPES.IChurch,
    dispatch: Dispatch<any>
): Array<JSX.Element> => {
    let data: Array<any> = []
    let columns: Array<Column<any>> = []

    tableDonators.forEach((d) => {
        const donations = d.donations
        if (donations) {
            donations.forEach((donation) => {
                data.push(
                    {
                        donatorid: donation?.donatorid || null,
                        churchid: donation?.churchid || null,
                        name: (d?.firstname || '') + ' ' + (d?.lastname || ''),
                        email: d?.email || null,
                        date: donation?.date || null,
                        amount: donation?.amount || null,
                        type: donation?.type || null,
                        currency: donation?.currency || null,
                        category: donation?.category || null,
                        details: donation?.details || null
                    }
                )
            })

        }
    })

    columns = [
        { title: 'Name', field: 'name', editable: 'never' },
        { title: 'Email', field: 'email', editable: 'never' },
        { title: 'Date', field: 'date', filtering: true, type: 'datetime', editable: 'never' },
        { title: 'Amount', field: 'amount', type: 'numeric' },
        { title: 'Type', field: 'type', lookup: DONATION_TYPES.DonationTypes },
        { title: 'Currency', field: 'currency', editable: 'never' },
        { title: 'Category', field: 'category', lookup: church.donationcategories },
        { title: 'Details', field: 'details' },
    ]

    const materialTables = [
        <MaterialTable
            key={'donation-table'}
            title={(church?.name || 'No') + ' Donations'}
            columns={columns}
            data={data}
            options={{
                grouping: true,
                filtering: true,
                exportButton: true,
                exportAllData: true
            }}
            editable={{
                onRowUpdate: (newDonation, oldDonation) =>
                    new Promise((resolve) => {
                        setTimeout(() => {
                            resolve();
                            if (oldDonation) {
                                if (
                                    newDonation.category !== oldDonation.category
                                    || newDonation.type !== oldDonation.type
                                ) {
                                    const deleteDonation: DONATION_TYPES.IDonation = {
                                        donatorid: parseInt(oldDonation?.donatorid),
                                        churchid: church?.id,
                                        date: oldDonation?.date,
                                        amount: parseFloat(oldDonation?.amount) * -1,
                                        type: oldDonation?.type,
                                        category: oldDonation?.category,
                                        currency: oldDonation?.currency,
                                        details: oldDonation?.details,
                                    }

                                    dispatch(updateDonationACT(deleteDonation))
                                    const donation: DONATION_TYPES.IDonation = {
                                        donatorid: parseInt(newDonation?.donatorid),
                                        churchid: church?.id,
                                        date: newDonation?.date,
                                        amount: parseFloat(newDonation?.amount),
                                        type: newDonation?.type,
                                        category: newDonation?.category,
                                        currency: newDonation?.currency,
                                        details: newDonation?.details,
                                    }
                                    dispatch(updateDonationACT(donation))
                                } else {
                                    const donation: DONATION_TYPES.IDonation = {
                                        donatorid: parseInt(newDonation?.donatorid),
                                        churchid: church?.id,
                                        date: newDonation?.date,
                                        amount: parseFloat(newDonation?.amount) - parseFloat(oldDonation?.amount),
                                        type: newDonation?.type,
                                        category: newDonation?.category,
                                        currency: newDonation?.currency,
                                        details: newDonation?.details,
                                    }
                                    dispatch(updateDonationACT(donation))
                                }

                            }
                        }, 600);
                    }),
                onRowDelete: oldDonation =>
                    new Promise((resolve, reject) => {
                        setTimeout(() => {
                            resolve();
                            const donation: DONATION_TYPES.IDonation = {
                                donatorid: parseInt(oldDonation?.donatorid),
                                churchid: church?.id,
                                date: oldDonation?.date,
                                amount: parseFloat(oldDonation?.amount) * -1,
                                type: oldDonation?.type,
                                category: oldDonation?.category,
                                currency: oldDonation?.currency,
                                details: oldDonation?.details,
                            }
                            dispatch(updateDonationACT(donation))
                        }, 1000);
                    })
            }}
        />
    ]
    return materialTables
}

const prepareChurchDonationReportTables = (
    donationReport: CHURCH_TYPES.IChurchDonationReport,
    church: CHURCH_TYPES.IChurch,
    donators: CHURCH_TYPES.IDonators
): Array<JSX.Element> => {
    let data: Array<any> = []
    let columns: Array<StaticDynamicColumns> = []

    let columnMap: { [index: string]: StaticDynamicColumns } = {
        'name': { title: 'Name', field: 'name', editable: 'never', static: true },
        'email': { title: 'Email', field: 'email', editable: 'never', static: true },
        'category': { title: 'Category', field: 'category', editable: 'never', static: true },
        'date': { title: 'Date', field: 'date', filtering: true, type: 'datetime', editable: 'never', static: true }
    }
    let totalsColumnMap: { [index: string]: StaticDynamicColumns } = {
        'category': { title: 'Category', field: 'category', editable: 'never', static: true }
    }
    let grandTotalsColumnMap: { [index: string]: StaticDynamicColumns } = {}

    Object.entries(donationReport.donations || {}).forEach(entry => {
        const donatorid: number = parseInt(entry[0])
        const donationReporttRow: { [columnName: string]: Array<DONATION_TYPES.IDonation> } = entry[1]

        Object.entries(donationReporttRow).forEach(entry => {
            const columnName: string = entry[0]
            let donations: Array<DONATION_TYPES.IDonation> = entry[1]

            columnMap[columnName] = {
                title: 'year :' + columnName.split('_')[0] + ' ' + donationReport.sumfilter.timeperiod + ' : ' + columnName.split('_')[1],
                field: columnName,
                editable: 'never',
                static: false
            }
            totalsColumnMap[columnName] = {
                title: 'year :' + columnName.split('_')[0] + ' ' + donationReport.sumfilter.timeperiod + ' : ' + columnName.split('_')[1],
                field: columnName,
                editable: 'never',
                static: false
            }
            grandTotalsColumnMap[columnName] = {
                title: 'year :' + columnName.split('_')[0] + ' ' + donationReport.sumfilter.timeperiod + ' : ' + columnName.split('_')[1],
                field: columnName,
                editable: 'never',
                static: false
            }

            donations.forEach(donation => {
                data.push(
                    {
                        name: (donators[donatorid].firstname || '') + ' ' + (donators[donatorid]?.lastname || ''),
                        email: donators[donatorid]?.email || null,
                        category: donation?.category || null,
                        date: donation?.date || null,
                        [columnName]: donation?.amount || null
                    }
                )
            })

        })
    })

    let totalData: Array<any> = []
    let grandTotalDataMap: { [columnName: string]: number } = {}
    let grandTotalData: Array<any> = []

    Object.entries(donationReport.donationssum || {}).forEach(entry => {
        const category: string = entry[0]
        const donationTotalReportRow: { [columnName: string]: number } = entry[1]
        totalData.push(
            {
                category: category,
                ...donationTotalReportRow
            }
        )
        // Add column values vertically to get grand totals
        Object.entries(donationTotalReportRow).forEach(entry => {
            const columnName: string = entry[0]
            const amount: number = entry[1]

            if (!grandTotalDataMap[columnName]) {
                grandTotalDataMap[columnName] = amount
            } else {
                grandTotalDataMap[columnName] += amount
            }

        })
    })

    grandTotalData.push({ ...grandTotalDataMap })

    columns = Object.values(columnMap)
    let totalColumns = Object.values(totalsColumnMap)
    totalColumns.push(
        { title: 'Total', field: 'total', editable: 'never', static: true }
    )
    let grandTotalColumns = Object.values(grandTotalsColumnMap)
    grandTotalColumns.push(
        { title: 'Grand Total', field: 'total', editable: 'never', static: true }
    )

    const materialTables = [
        <MaterialTable
            key={'church-donation-report-table'}
            title={(church?.name || 'No') + ' Donation Report'}
            columns={columns}
            data={data}
            options={{

                grouping: true,
                filtering: true,
                exportButton: true,
                exportAllData: true
            }}
        />,
        <MaterialTable
            key={'church-donation-report-totals-table'}
            title={(church?.name || 'No') + ' Totals'}
            columns={totalColumns}
            data={totalData}
            options={{
                grouping: true,
                filtering: true,
                exportButton: true,
                exportAllData: true

            }}
        />,
        <MaterialTable
            key={'church-donation-report-grandtotals-table'}
            title={(church?.name || 'No') + ' Grand Totals'}
            columns={grandTotalColumns}
            data={grandTotalData}
            options={{
                grouping: true,
                filtering: true,
                exportButton: true,
                exportAllData: true
            }}
        />
    ]
    return materialTables
}
// TODO fix loading problem with large number of material tables
const prepareDonationStatementReportTables = (
    statementReport: CHURCH_TYPES.IDonationStatementReport,
    church: CHURCH_TYPES.IChurch,
    donators: CHURCH_TYPES.IDonators
): Array<JSX.Element> => {
    let materialTables: Array<JSX.Element> = []

    let columns: Array<Column<any>> = [
        { title: 'Date', field: 'date', filtering: true, type: 'datetime', editable: 'never' },
        { title: 'Amount', field: 'amount', type: 'numeric' },
        { title: 'Type', field: 'type', lookup: DONATION_TYPES.DonationTypes },
        { title: 'Currency', field: 'currency', editable: 'never' },
        { title: 'Category', field: 'category', lookup: church.donationcategories },
        { title: 'Details', field: 'details' },
    ]
    let totalColumns: Array<Column<any>> = [
        { title: 'Category', field: 'category', lookup: church.donationcategories },
        { title: 'Total', field: 'total', type: 'numeric' },
    ]
    let grandTotalColumns: Array<Column<any>> = [
        { title: 'Grand Total', field: 'total', type: 'numeric' },
    ]

    Object.entries(statementReport.donations || {}).forEach(entry => {
        const donatorid: number = parseInt(entry[0])
        const donator = donators[donatorid]

        let data: Array<DONATION_TYPES.IDonation> = []
        let totalData: Array<{ total: number, category: string }> = []
        let grandTotalData: Array<{ total: number }> = []

        const statementReportRow: { [key: string]: Array<DONATION_TYPES.IDonation> } = entry[1]
        Object.values(statementReportRow).forEach(donations => {
            donations.forEach(donation => {
                data.push(
                    {
                        date: donation?.date || undefined,
                        amount: donation?.amount || undefined,
                        type: donation?.type || undefined,
                        currency: donation?.currency || undefined,
                        category: donation?.category || undefined,
                        details: donation?.details || undefined
                    }
                )
            })
        })

        const statementTotalsPerCategory: { [category: string]: { [donatorIDdateCategoryKey: string]: number } } | undefined = statementReport.donationssum
        Object.entries(statementTotalsPerCategory || {}).forEach(entry => {
            const category = entry[0]
            const statementAmounts = entry[1]
            const total = statementAmounts[donatorid + "_total"] || 0

            totalData.push(
                {
                    total: total,
                    category: category,
                }
            )

        })

        const grandTotal = totalData.reduce((acc, curr) => acc += curr.total, 0)
        grandTotalData.push({
            total: grandTotal
        })


        materialTables = [
            ...materialTables,
            <MaterialTable
                key={'donation-statement-report-table'}
                title={' Donation Statement Report for ' + (donator?.firstname + ' ' + donator?.lastname)}
                columns={columns}
                data={data}
                options={{

                    grouping: true,
                    filtering: true,
                    exportButton: true,
                    exportAllData: true
                }}
            />,
            <MaterialTable
                key={'donation-statement-report-totals-table'}
                title={' Totals for ' + (donator?.firstname + ' ' + donator?.lastname)}
                columns={totalColumns}
                data={totalData}
                options={{
                    grouping: true,
                    filtering: true,
                    exportButton: true,
                    exportAllData: true

                }}
            />,
            <MaterialTable
                key={'donation-statement-report-grandtotals-table'}
                title={'Grand Total for ' + (donator?.firstname + ' ' + donator?.lastname)}
                columns={grandTotalColumns}
                data={grandTotalData}
                options={{
                    grouping: true,
                    filtering: true,
                    exportButton: true,
                    exportAllData: true
                }}
            />
        ]
    })

    return materialTables
}

const exportAllTablesToPDF = (church: CHURCH_TYPES.IChurch, materialTables: Array<JSX.Element>) => {
    let doc = new jsPDF();
    materialTables.forEach((materialTable) => {
        const props = (materialTable.props as MaterialTableProps<any>)

        const columns = (materialTable.props as MaterialTableProps<any>).columns as StaticDynamicColumns[]
        const staticColumnNames: Array<ColumnInput> = columns.filter(column => column.static).map(column => { return { title: column.title as string, key: column.field as string } })
        const staticColumnWidths = columns.filter(column => column.static).map(column => { return { [column.field as string]: { cellWidth: 25 } } }).reduce((acc, curr) => {
            const entry = Object.entries(curr)[0]
            const key = entry[0]
            const value = entry[1]
            acc[key] = value
            return acc
        }, {})
        const dynamicColumnNames: Array<ColumnInput> = columns.filter(column => !column.static).map(column => { return { title: column.title as string, key: column.field as string } })
        const dynamicSplitColumnNames = dynamicColumnNames.length !== 0 ? splitIntoNLengthSubArrays(dynamicColumnNames, 7 - staticColumnNames.length) : null

        const data = props.data as Array<any>
        const body: Array<{ [key: string]: string }> = data.map(dataObject => {
            return columns.map(name => {
                const cellInput: CellInput = dataObject[name.field || ""] as string
                return { [name.field as string]: cellInput }
            }).reduce((acc, curr) => {
                Object.assign(acc, curr)
                return acc
            }, {})
        })

        const yPos = (doc as any).lastAutoTable.finalY
        doc.text(props.title, 14, yPos + 10 || 20)

        if (dynamicSplitColumnNames) {
            dynamicSplitColumnNames.forEach(dynamicColumnNames => {
                const staticKeySet = staticColumnNames.map((cInput) => (cInput as { title: string, key: string }).key)
                const dynamicKeySet = dynamicColumnNames.map((cInput) => (cInput as { title: string, key: string }).key)
                const keySet = [
                    ...staticKeySet,
                    ...dynamicKeySet
                ]

                const rows = intersectArrayOfObjectSetArraysWithKeySet(
                    body,
                    keySet
                )
                const partialBody = rows.filter(row => {
                    let dynamicColumnValues: Array<any> = []
                    Object.entries(row).forEach(entry => {
                        const key = entry[0]
                        const value = entry[1]
                        if (dynamicKeySet.includes(key)) {
                            dynamicColumnValues.push(value)
                        }
                    })
                    return !dynamicColumnValues.every(v => v === undefined)
                })

                const autoTableColumns = [...staticColumnNames, ...dynamicColumnNames]
                const autoTableBody = Object.values(partialBody)
                const yPos = (doc as any).lastAutoTable.finalY
                autoTable(
                    doc,
                    {
                        columnStyles: { ...staticColumnWidths },
                        columns: autoTableColumns,
                        body: autoTableBody,
                        startY: yPos + 15 || 25
                    }
                )
            })
        } else {
            autoTable(doc, { columns: staticColumnNames, body: body, startY: yPos + 15 || 25 })
        }
    })
    doc.save(church?.name + "_report.pdf")
}

const splitIntoNLengthSubArrays = <T,>(a: Array<T>, subArrayLength: number): Array<Array<T>> => {
    let len: number = a.length
    if (subArrayLength >= len) { return [a]; }

    let out: Array<any> = []
    let i: number = 0

    while (i < len) {
        out.push(a.slice(i, i += subArrayLength));
    }

    return out;

}

const intersectArrayOfObjectSetArraysWithKeySet = (objectSetArray: Array<{ [key: string]: any }>, keySet: Array<string>): Array<{ [key: string]: any }> => {
    const newObjectSetArray: Array<{ [key: string]: any }> = []

    objectSetArray.forEach(v => {
        let newObjectSet: { [index: string]: any } = {}

        Object.entries(v).forEach(entry => {
            const key = entry[0]
            const value = entry[1]

            if (keySet.includes(key)) {
                newObjectSet[key] = value
            }
        })
        newObjectSetArray.push(newObjectSet)
    })

    return newObjectSetArray
}