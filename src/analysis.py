import pandas as pd
import matplotlib.pyplot as plt
import numpy as np
from datetime import datetime


font = {
        # 'family': 'serif',
        # 'color':  'darkred',
        'weight': 'normal',
        'size': 12,
        }


def plotMetrics(df, figName, savePath=None):
    fig, (ax1, ax2, ax3, ax4) = plt.subplots(4, 1, figsize=(16, 8), sharex='col')
    fig.suptitle(figName, fontsize=14)

    x = df.index
    p0 = df['MarkPrice'][0]
    e0 = df['Equity'][0]

    ax1.set_title('Price')
    ax1.plot(x, np.array(df['MarkPrice']), linewidth=1.3)
    lbl = ax1.set_ylabel('$', labelpad=10)
    lbl.set_rotation(0)
    ax1b = ax1.twinx()
    ax1b.plot(x, (df['MarkPrice'] - p0) / p0 * 100, linewidth=0)
    lbl = ax1b.set_ylabel('%', labelpad=10)
    lbl.set_rotation(0)

    ax2.set_title('Equity')
    ax2.plot(x, df['Equity'], linewidth=1)
    ax2.fill_between(x, df['Equity'][0], df['Equity'], alpha=0.5)
    lbl = ax2.set_ylabel('$', labelpad=10)
    lbl.set_rotation(0)
    ax2.ticklabel_format(axis='y', useOffset=False)
    ax2b = ax2.twinx()
    ax2b.plot(x, (df['Equity'] - e0) / e0 * 100, linewidth=0)
    lbl = ax2b.set_ylabel('%', labelpad=10)
    lbl.set_rotation(0)

    ax3.set_title('PNL-L')
    ax3.plot(x, df['PNL-L'], linewidth=1, color='red')
    ax3.fill_between(x, df['PNL-L'][0], df['PNL-L'], alpha=0.5, color='red')
    lbl = ax3.set_ylabel('$', labelpad=10)
    lbl.set_rotation(0)
    ax3.ticklabel_format(axis='y', useOffset=False)
    ax3b = ax3.twinx()
    ax3b.plot(x, df['PNL-L'] / df['Equity'] * 100, linewidth=0)
    lbl = ax3b.set_ylabel('%', labelpad=10)
    lbl.set_rotation(0)

    ax4.set_title('PNL-S')
    ax4.plot(x, df['PNL-S'], linewidth=1, color='red')
    ax4.fill_between(x, df['PNL-S'][0], df['PNL-S'], alpha=0.5, color='red')
    lbl = ax4.set_ylabel('$', labelpad=10)
    lbl.set_rotation(0)
    ax4.ticklabel_format(axis='y', useOffset=False)
    ax4b = ax4.twinx()
    ax4b.plot(x, df['PNL-S'] / df['Equity'] * 100, linewidth=0)
    lbl = ax4b.set_ylabel('%', labelpad=10)
    lbl.set_rotation(0)
  
    fig.tight_layout()
    if savePath is not None:
        fig.savefig(savePath)


def plotDistributions(df, side, savePath=None):
    profit = 'NetProfit-' + side
    gridReached = 'GridReached-'+ side

    df[profit] = df['GrossProfit-L']
    fig, (ax1, ax2) = plt.subplots(1, 2, figsize=(16, 6))
    # fig.suptitle('Distribution of take profits', fontsize=14)

    # get grid reached for each profit
    profitGridReached = []
    for i in range(len(df.index)):
        if df[profit][i] != 0:
            profitGridReached.append(df[gridReached][i-1]) # TODO change i-1

    data = np.array(profitGridReached)
    bins = np.arange(-1, data.max() + 0.5) - 0.5
    ax1.hist(data, bins, rwidth=0.7, density=True)
    ax1.set_xticks(bins + 0.5)
    ax1.set_xlim([0.1-1, data.max() + 0.9])
    ax1.set_title('Distribution of TP grids', fontdict=font)
    ax1.set_ylabel('Frequency', fontdict=font)
    ax1.set_xlabel('Grid Number', fontdict=font)

    # compute profit contribution of each grid
    gridProfitDict = {}
    for i in range(len(df.index)):
        netProfit = df[profit][i]

        if netProfit != 0:
            grid = int(df[gridReached][i-1]) # TODO change i-1
            if grid in gridProfitDict:
                gridProfitDict[grid] += netProfit
            else:
                gridProfitDict[grid] = netProfit

    # check the sum of grid profits is equal to the total profit
    totalNetProfit = df[profit].sum()
    if abs(totalNetProfit - sum(list(gridProfitDict.values()))) > 0.0001:
        print(totalNetProfit)
        print(sum(list(gridProfitDict.values())))
        raise Exception("The sum of the grid profits is not equal to the total profits")

    nGrids = int(max(list(gridProfitDict.keys())))
    data = []
    for i in range (nGrids):
        if (i+1) in gridProfitDict:
            data.append(gridProfitDict[i+1])
            # data.append(gridProfitDict[i+1] / abs(df[profit]).sum() * 100)
        else:
            gridProfitDict[i] = 0
            data.append(0)

    # distribution of profits
    ax2.bar(range(0,nGrids), data)
    ax2.set_xticks(bins + 0.5)
    ax2.set_xlim([0.1-1, nGrids + 0.9])
    ax2.set_title('Distribution of profit per grid', fontdict=font)
    # ax2.set_ylabel('Profit %', fontdict=font)
    ax2.set_xlabel('Grid Number', fontdict=font)

    ax2.set_ylabel('Profit $', labelpad=10, fontdict=font)
    ax2b = ax2.twinx()
    ax2b.bar(range(0,nGrids), data / abs(df[profit]).sum() * 100)
    ax2b.set_ylabel('Profit %', labelpad=10, fontdict=font)

    fig.tight_layout()
    fig.subplots_adjust(wspace=0.15)
    if savePath is not None:
        fig.savefig(savePath)


if __name__ == '__main__':
    folder = "results/"
    file = "AntiMartingala LONG GO 7, GS 0.30, SF 1.50, OS 1.00, OF 2.00, TS 0.30, SL 0.30.csv"
    df = pd.read_csv(folder + file)
    df.index = [datetime.fromtimestamp(x) for x in df['Timestamp']]

    plotMetrics(df, file)
    plotDistributions(df, 'L')
    plt.show()