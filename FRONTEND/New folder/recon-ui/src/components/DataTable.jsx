import { Table, TableBody, TableCell, TableContainer, TableHead, TableRow, Paper } from '@mui/material';

export default function DataTable({ columns, rows, onRowClick }) {
  return (
    <TableContainer component={Paper} variant="outlined">
      <Table>
        <TableHead>
          <TableRow>
            {columns.map((c) => (
              <TableCell key={c.key} align={c.align}>{c.label}</TableCell>
            ))}
          </TableRow>
        </TableHead>
        <TableBody>
          {rows.map((r, i) => (
            <TableRow
              key={r.id ?? i}
              hover={!!onRowClick}
              onClick={onRowClick ? () => onRowClick(r) : undefined}
              sx={{ cursor: onRowClick ? 'pointer' : 'default', '&:last-child td': { border: 0 } }}
            >
              {columns.map((c) => (
                <TableCell key={c.key} align={c.align} sx={{ py: 1.5, fontWeight: c.strong ? 500 : 400 }}>
                  {c.render ? c.render(r) : r[c.key]}
                </TableCell>
              ))}
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </TableContainer>
  );
}
