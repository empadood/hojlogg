import React from 'react';
import { View, Text, TouchableOpacity, StyleSheet } from 'react-native';
import type { Log } from '../api/client';

interface Props {
  log: Log;
  onPress: () => void;
}

export function LogItem({ log, onPress }: Props) {
  const date = new Date(log.created_at).toLocaleDateString(undefined, {
    day: '2-digit',
    month: 'short',
    year: 'numeric',
  });

  return (
    <TouchableOpacity style={styles.card} onPress={onPress} activeOpacity={0.75}>
      <View style={styles.row}>
        <Text style={styles.odometer}>
          {log.odometer_km.toLocaleString()} km
        </Text>
        <Text style={styles.date}>{date}</Text>
      </View>
      {log.fuel_level != null && (
        <View style={styles.row}>
          <Text style={styles.label}>Fuel</Text>
          <Text style={styles.value}>{log.fuel_level.toFixed(0)} %</Text>
        </View>
      )}
      {log.notes ? <Text style={styles.notes} numberOfLines={2}>{log.notes}</Text> : null}
      {log.parsed_by_ocr && (
        <View style={styles.ocrBadge}>
          <Text style={styles.ocrText}>OCR</Text>
        </View>
      )}
    </TouchableOpacity>
  );
}

const CARD_BG = '#1e1e2e';
const ACCENT = '#f97316';

const styles = StyleSheet.create({
  card: {
    backgroundColor: CARD_BG,
    borderRadius: 12,
    padding: 16,
    marginHorizontal: 16,
    marginVertical: 6,
    shadowColor: '#000',
    shadowOpacity: 0.3,
    shadowRadius: 6,
    shadowOffset: { width: 0, height: 3 },
    elevation: 4,
  },
  row: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 4,
  },
  odometer: {
    fontSize: 20,
    fontWeight: '700',
    color: '#fff',
  },
  date: {
    fontSize: 13,
    color: '#aaa',
  },
  label: {
    fontSize: 13,
    color: '#aaa',
  },
  value: {
    fontSize: 13,
    color: '#fff',
  },
  notes: {
    fontSize: 13,
    color: '#ccc',
    marginTop: 4,
  },
  ocrBadge: {
    alignSelf: 'flex-start',
    backgroundColor: ACCENT,
    borderRadius: 4,
    paddingHorizontal: 6,
    paddingVertical: 2,
    marginTop: 6,
  },
  ocrText: {
    fontSize: 10,
    fontWeight: '700',
    color: '#fff',
  },
});
