import React, { useCallback, useState } from 'react';
import {
  View,
  Text,
  ScrollView,
  StyleSheet,
  Image,
  ActivityIndicator,
  Alert,
  TouchableOpacity,
} from 'react-native';
import { useFocusEffect } from '@react-navigation/native';
import type { NativeStackScreenProps } from '@react-navigation/native-stack';
import { logsApi, type Log, BASE_URL } from '../api/client';
import type { RootStackParamList } from '../../App';

type Props = NativeStackScreenProps<RootStackParamList, 'LogDetail'>;

export function LogDetailScreen({ navigation, route }: Props) {
  const { id } = route.params;
  const [log, setLog] = useState<Log | null>(null);
  const [loading, setLoading] = useState(true);

  useFocusEffect(
    useCallback(() => {
      logsApi
        .get(id)
        .then(setLog)
        .catch(() => Alert.alert('Error', 'Could not load log.'))
        .finally(() => setLoading(false));
    }, [id])
  );

  const handleDelete = () => {
    Alert.alert('Delete', 'Delete this log?', [
      { text: 'Cancel', style: 'cancel' },
      {
        text: 'Delete',
        style: 'destructive',
        onPress: async () => {
          await logsApi.delete(id).catch(() => null);
          navigation.goBack();
        },
      },
    ]);
  };

  if (loading) {
    return (
      <View style={styles.center}>
        <ActivityIndicator color={ACCENT} size="large" />
      </View>
    );
  }

  if (!log) {
    return (
      <View style={styles.center}>
        <Text style={styles.error}>Log not found.</Text>
      </View>
    );
  }

  const date = new Date(log.created_at).toLocaleString();

  return (
    <ScrollView style={styles.container} contentContainerStyle={styles.content}>
      {log.image_path ? (
        <Image
          source={{ uri: `${BASE_URL}/uploads/${log.image_path.split('/').pop()}` }}
          style={styles.image}
          resizeMode="cover"
        />
      ) : (
        <TouchableOpacity
          style={styles.noImage}
          onPress={() => navigation.navigate('Camera', { logId: log.id })}
        >
          <Text style={styles.noImageIcon}>📷</Text>
          <Text style={styles.noImageText}>Tap to add a dashboard photo</Text>
        </TouchableOpacity>
      )}

      <View style={styles.section}>
        <Row label="Odometer" value={`${log.odometer_km.toLocaleString()} km`} highlight />
        {log.fuel_level != null && (
          <Row label="Fuel level" value={`${log.fuel_level.toFixed(0)} %`} />
        )}
        <Row label="Date" value={date} />
        {log.parsed_by_ocr && <Row label="Source" value="OCR scan" />}
        {log.notes ? <Row label="Notes" value={log.notes} /> : null}
      </View>

      <TouchableOpacity style={styles.deleteBtn} onPress={handleDelete}>
        <Text style={styles.deleteBtnText}>Delete Log</Text>
      </TouchableOpacity>
    </ScrollView>
  );
}

function Row({ label, value, highlight }: { label: string; value: string; highlight?: boolean }) {
  return (
    <View style={styles.row}>
      <Text style={styles.rowLabel}>{label}</Text>
      <Text style={[styles.rowValue, highlight && styles.rowValueHighlight]}>{value}</Text>
    </View>
  );
}

const BG = '#0f0f1a';
const CARD = '#1e1e2e';
const ACCENT = '#f97316';

const styles = StyleSheet.create({
  container: { flex: 1, backgroundColor: BG },
  content: { paddingBottom: 40 },
  center: { flex: 1, backgroundColor: BG, alignItems: 'center', justifyContent: 'center' },
  error: { color: '#f55', fontSize: 16 },
  image: { width: '100%', height: 220 },
  noImage: {
    height: 180,
    backgroundColor: CARD,
    alignItems: 'center',
    justifyContent: 'center',
    margin: 16,
    borderRadius: 12,
  },
  noImageIcon: { fontSize: 40 },
  noImageText: { color: '#888', marginTop: 8 },
  section: {
    backgroundColor: CARD,
    borderRadius: 12,
    marginHorizontal: 16,
    marginTop: 8,
    padding: 16,
  },
  row: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    paddingVertical: 10,
    borderBottomWidth: 1,
    borderBottomColor: '#2a2a3e',
  },
  rowLabel: { color: '#aaa', fontSize: 14 },
  rowValue: { color: '#fff', fontSize: 14 },
  rowValueHighlight: { color: ACCENT, fontWeight: '700', fontSize: 18 },
  deleteBtn: {
    margin: 16,
    marginTop: 24,
    padding: 16,
    borderRadius: 10,
    borderWidth: 1,
    borderColor: '#f55',
    alignItems: 'center',
  },
  deleteBtnText: { color: '#f55', fontWeight: '600', fontSize: 15 },
});
