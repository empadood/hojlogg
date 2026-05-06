import React, { useState } from 'react';
import {
  View,
  Text,
  TextInput,
  TouchableOpacity,
  StyleSheet,
  ScrollView,
  Alert,
  Image,
  ActivityIndicator,
  KeyboardAvoidingView,
  Platform,
} from 'react-native';
import type { NativeStackScreenProps } from '@react-navigation/native-stack';
import { logsApi } from '../api/client';
import type { RootStackParamList } from '../../App';

type Props = NativeStackScreenProps<RootStackParamList, 'NewLog'>;

export function NewLogScreen({ navigation, route }: Props) {
  const { imageUri, logId } = route.params ?? {};

  const [odometer, setOdometer] = useState('');
  const [fuelLevel, setFuelLevel] = useState('');
  const [notes, setNotes] = useState('');
  const [saving, setSaving] = useState(false);

  const handleSave = async () => {
    const odometerNum = parseFloat(odometer);
    if (isNaN(odometerNum) || odometerNum < 0) {
      Alert.alert('Validation', 'Please enter a valid odometer reading (km).');
      return;
    }

    const fuelNum = fuelLevel.trim() !== '' ? parseFloat(fuelLevel) : undefined;
    if (fuelNum !== undefined && (isNaN(fuelNum) || fuelNum < 0 || fuelNum > 100)) {
      Alert.alert('Validation', 'Fuel level must be between 0 and 100.');
      return;
    }

    setSaving(true);
    try {
      const log = await logsApi.create({
        odometer_km: odometerNum,
        fuel_level: fuelNum ?? null,
        notes: notes.trim(),
      });

      if (imageUri) {
        try {
          const result = await logsApi.uploadImage(log.id, imageUri);
          // If OCR produced values, let the user know.
          if (result.ocr.odometer_km != null) {
            Alert.alert(
              'OCR Result',
              `Dashboard analysed!\nOdometer: ${result.ocr.odometer_km.toLocaleString()} km\nFuel: ${result.ocr.fuel_level?.toFixed(0) ?? '?'} %\nConfidence: ${(result.ocr.confidence * 100).toFixed(0)} %`
            );
          }
        } catch {
          // Image upload failure is non-fatal.
          Alert.alert('Warning', 'Log saved but image upload failed.');
        }
      }

      navigation.goBack();
    } catch (err) {
      Alert.alert('Error', 'Could not save log. Is the backend running?');
    } finally {
      setSaving(false);
    }
  };

  return (
    <KeyboardAvoidingView
      style={styles.root}
      behavior={Platform.OS === 'ios' ? 'padding' : undefined}
    >
      <ScrollView contentContainerStyle={styles.container}>
        {/* Dashboard image preview */}
        <TouchableOpacity
          style={styles.imagePlaceholder}
          onPress={() => navigation.navigate('Camera', { logId })}
          activeOpacity={0.8}
        >
          {imageUri ? (
            <Image source={{ uri: imageUri }} style={styles.image} resizeMode="cover" />
          ) : (
            <>
              <Text style={styles.cameraIcon}>📷</Text>
              <Text style={styles.cameraHint}>Tap to scan dashboard</Text>
            </>
          )}
        </TouchableOpacity>

        <Text style={styles.label}>Odometer (km) *</Text>
        <TextInput
          style={styles.input}
          placeholder="e.g. 12345"
          placeholderTextColor="#555"
          keyboardType="decimal-pad"
          value={odometer}
          onChangeText={setOdometer}
        />

        <Text style={styles.label}>Fuel level (%)</Text>
        <TextInput
          style={styles.input}
          placeholder="0 – 100"
          placeholderTextColor="#555"
          keyboardType="decimal-pad"
          value={fuelLevel}
          onChangeText={setFuelLevel}
        />

        <Text style={styles.label}>Notes</Text>
        <TextInput
          style={[styles.input, styles.textarea]}
          placeholder="Optional notes…"
          placeholderTextColor="#555"
          multiline
          numberOfLines={4}
          value={notes}
          onChangeText={setNotes}
        />

        <TouchableOpacity
          style={[styles.saveBtn, saving && styles.saveBtnDisabled]}
          onPress={handleSave}
          disabled={saving}
        >
          {saving ? (
            <ActivityIndicator color="#fff" />
          ) : (
            <Text style={styles.saveBtnText}>Save Log</Text>
          )}
        </TouchableOpacity>
      </ScrollView>
    </KeyboardAvoidingView>
  );
}

const BG = '#0f0f1a';
const CARD = '#1e1e2e';
const ACCENT = '#f97316';

const styles = StyleSheet.create({
  root: { flex: 1, backgroundColor: BG },
  container: { padding: 20, paddingBottom: 40 },
  imagePlaceholder: {
    height: 200,
    backgroundColor: CARD,
    borderRadius: 12,
    alignItems: 'center',
    justifyContent: 'center',
    marginBottom: 24,
    overflow: 'hidden',
  },
  image: { width: '100%', height: '100%' },
  cameraIcon: { fontSize: 48 },
  cameraHint: { color: '#888', marginTop: 8 },
  label: { color: '#ccc', fontSize: 13, marginBottom: 6, marginTop: 12 },
  input: {
    backgroundColor: CARD,
    color: '#fff',
    borderRadius: 8,
    padding: 14,
    fontSize: 16,
    borderWidth: 1,
    borderColor: '#2a2a3e',
  },
  textarea: { height: 100, textAlignVertical: 'top' },
  saveBtn: {
    backgroundColor: ACCENT,
    padding: 16,
    borderRadius: 10,
    alignItems: 'center',
    marginTop: 28,
  },
  saveBtnDisabled: { opacity: 0.6 },
  saveBtnText: { color: '#fff', fontWeight: '700', fontSize: 16 },
});
