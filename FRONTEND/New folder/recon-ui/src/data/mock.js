export const docs = [
  { id: 1, name: 'INV-2043_Keells.pdf', source: 'WhatsApp', type: 'Supplier invoice', company: 'Keells Super (Pvt) Ltd', received: 'Oct 27 · 10:42', confidence: '98%', status: 'Processed' },
  { id: 2, name: 'PO-8812_Hayleys.pdf', source: 'Gmail', type: 'Purchase order', company: 'Hayleys PLC', received: 'Oct 27 · 10:36', confidence: '96%', status: 'Processing' },
  { id: 3, name: 'BT-2210_Slip.pdf', source: 'Manual upload', type: 'Bank transfer', company: 'Dialog Axiata PLC', received: 'Oct 27 · 10:21', confidence: '74%', status: 'Needs review' },
  { id: 4, name: 'INV-2039_Cargills.pdf', source: 'Gmail', type: 'Supplier invoice', company: 'Cargills (Ceylon) PLC', received: 'Oct 27 · 09:58', confidence: '99%', status: 'Processed' },
  { id: 5, name: 'scan_0091.pdf', source: 'WhatsApp', type: 'Unclassified', company: '—', received: 'Oct 27 · 09:40', confidence: '31%', status: 'Failed' },
  { id: 6, name: 'INV-2037_Abans.pdf', source: 'WhatsApp', type: 'Supplier invoice', company: 'Abans PLC', received: 'Oct 27 · 09:12', confidence: '97%', status: 'Processed' },
  { id: 7, name: 'Remittance_1104.pdf', source: 'Gmail', type: 'Bank transfer', company: 'Sampath Bank', received: 'Oct 26 · 16:47', confidence: '95%', status: 'Processed' },
  { id: 8, name: 'PO-8809_LOLC.pdf', source: 'Gmail', type: 'Purchase order', company: 'LOLC Finance', received: 'Oct 26 · 15:30', confidence: '93%', status: 'Processing' },
  { id: 9, name: 'INV-2036_Softlogic.pdf', source: 'Manual upload', type: 'Supplier invoice', company: 'Softlogic PLC', received: 'Oct 26 · 14:05', confidence: '98%', status: 'Processed' },
];

export const whatsappFeed = [
  { name: 'INV-2043_Keells.pdf', meta: 'Keells Super · 2m ago', status: 'Processed' },
  { name: 'PO-8812_Hayleys.pdf', meta: 'Hayleys PLC · 6m ago', status: 'Processing' },
  { name: 'slip_20261025.pdf', meta: 'Dialog Axiata · 14m ago', status: 'Needs review' },
  { name: 'INV-2037_Abans.pdf', meta: 'Abans PLC · 31m ago', status: 'Processed' },
  { name: 'scan_0091.pdf', meta: 'Unknown sender · 1h ago', status: 'Failed' },
];

export const gmailFeed = [
  { name: 'INV-2039_Cargills.pdf', meta: 'Cargills · 1m ago', status: 'Processed' },
  { name: 'Statement_Oct.pdf', meta: 'Commercial Bank · 9m ago', status: 'Needs review' },
  { name: 'INV-2036_Softlogic.pdf', meta: 'Softlogic · 22m ago', status: 'Processed' },
  { name: 'PO-8809_LOLC.pdf', meta: 'LOLC Finance · 48m ago', status: 'Processing' },
  { name: 'Remittance_1104.pdf', meta: 'Sampath Bank · 1h ago', status: 'Processed' },
];

export const exceptions = [
  { id: 'EX-0412', doc: 'INV-2043_Keells', company: 'Keells Super (Pvt) Ltd', issue: 'Amount mismatch', variance: '+12,500.00', age: '2 h', status: 'Open' },
  { id: 'EX-0411', doc: 'INV-2041_Hayleys', company: 'Hayleys PLC', issue: 'Missing PO', variance: '—', age: '3 h', status: 'Open' },
  { id: 'EX-0409', doc: 'BT-2210_Slip', company: 'Dialog Axiata PLC', issue: 'Unmatched payment', variance: '-85,000.00', age: '5 h', status: 'In review' },
  { id: 'EX-0406', doc: 'INV-2038_Abans', company: 'Abans PLC', issue: 'Duplicate invoice', variance: '—', age: '7 h', status: 'In review' },
  { id: 'EX-0402', doc: 'INV-2035_Softlogic', company: 'Softlogic PLC', issue: 'Amount mismatch', variance: '+3,200.00', age: '1 d', status: 'Escalated' },
  { id: 'EX-0398', doc: 'scan_0091', company: 'Unknown', issue: 'Low OCR confidence', variance: '—', age: '1 d', status: 'Open' },
  { id: 'EX-0395', doc: 'INV-2031_Cargills', company: 'Cargills (Ceylon) PLC', issue: 'Missing PO', variance: '—', age: '2 d', status: 'Escalated' },
];

export const lineItems = [
  { item: 'Rice – Samba 25kg', poQty: 50, poPrice: '2,100.00', invQty: 50, invPrice: '2,100.00', variance: '0.00' },
  { item: 'Coconut oil 5L', poQty: 40, poPrice: '1,850.00', invQty: 40, invPrice: '1,850.00', variance: '0.00' },
  { item: 'Sugar 50kg', poQty: 25, poPrice: '11,000.00', invQty: 25, invPrice: '11,500.00', variance: '+12,500.00' },
  { item: 'Dhal – Red 10kg', poQty: 10, poPrice: '3,350.00', invQty: 10, invPrice: '3,350.00', variance: '0.00' },
];

export const ledger = [
  { id: '#12480', time: 'Oct 27 · 10:44', event: 'Settlement matched', invoice: 'INV-2043', amount: '500,000.00', hash: '9f3a…c21d', status: 'Verified' },
  { id: '#12479', time: 'Oct 27 · 10:41', event: 'Exception raised', invoice: 'INV-2043', amount: '12,500.00', hash: '41be…07a9', status: 'Verified' },
  { id: '#12478', time: 'Oct 27 · 10:12', event: '3-way match passed', invoice: 'INV-2039', amount: '1,240,000.00', hash: 'c8d0…5e12', status: 'Verified' },
  { id: '#12477', time: 'Oct 27 · 09:58', event: 'Invoice approved', invoice: 'INV-2037', amount: '318,750.00', hash: '7ab4…92f0', status: 'Verified' },
  { id: '#12476', time: 'Oct 27 · 09:30', event: 'Payment released', invoice: 'INV-2036', amount: '96,400.00', hash: 'e215…3bc8', status: 'Verified' },
  { id: '#12475', time: 'Oct 26 · 17:02', event: 'Exception raised', invoice: 'INV-2035', amount: '3,200.00', hash: '05f7…d46a', status: 'Verified' },
];

export const settlements = [
  { id: 'BT-2210', payer: 'Dialog Axiata PLC', amount: '500,000.00', invoice: 'INV-2043', date: 'Oct 25', status: 'Matched' },
  { id: 'BT-2207', payer: 'Cargills (Ceylon) PLC', amount: '1,240,000.00', invoice: 'INV-2039', date: 'Oct 24', status: 'Matched' },
  { id: 'BT-2204', payer: 'Sampath Bank', amount: '85,000.00', invoice: '—', date: 'Oct 24', status: 'Pending' },
  { id: 'BT-2201', payer: 'Hayleys PLC', amount: '210,000.00', invoice: 'INV-2033', date: 'Oct 23', status: 'Mismatch' },
];

export const users = [
  { id: 1, name: 'Admin User', email: 'admin@company.lk', role: 'Admin', last: 'Just now', status: 'Active' },
  { id: 2, name: 'Nimali Perera', email: 'nimali@company.lk', role: 'Reviewer', last: '12 min ago', status: 'Active' },
  { id: 3, name: 'Kasun Silva', email: 'kasun@company.lk', role: 'Reviewer', last: '2 h ago', status: 'Active' },
  { id: 4, name: 'Tharindu Fernando', email: 'tharindu@company.lk', role: 'Viewer', last: 'Yesterday', status: 'Active' },
  { id: 5, name: 'Ishara Jayasuriya', email: 'ishara@company.lk', role: 'Viewer', last: '—', status: 'Invited' },
  { id: 6, name: 'Dilan Wickramasinghe', email: 'dilan@company.lk', role: 'Reviewer', last: '3 weeks ago', status: 'Disabled' },
];
